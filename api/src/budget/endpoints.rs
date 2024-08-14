use actix_web::{post, get, web, HttpResponse, Responder};
use serde::{Serialize, Deserialize};
use sqlx::{MySqlPool, query_as};
use serde_json::json;
use crate::budget::helpers::{process_agency_data, AgencyBudget, make_map_data};


#[get("/agency")]
async fn get_agency(pool: web::Data<MySqlPool>) -> impl Responder {

    let result = query_as!(AgencyBudget, "Select agency, budget FROM agency_budget order by budget desc").fetch_all(pool.get_ref()).await;

    match result {
      Ok(rows) => {
        let (main_data, other_data, table_data) = process_agency_data(&rows);

        HttpResponse::Ok().json(json!({
            "main_data": main_data,
            "other_data": other_data,
            "table_data": table_data
        }))
    }
         Err(e) => {
            eprintln!("Failed to execute query: {}", e);
            HttpResponse::InternalServerError().json(json!({"error": "Failed to agency budget"}))
        }
    }
}

#[derive(Serialize)]
struct Country{
    country: Option<String>,
}

#[get("/get_countries")]
async fn get_countries(pool: web::Data<MySqlPool>) -> impl Responder {
    let result = query_as!(Country, "SELECT DISTINCT country FROM foreign_aid ORDER BY country").fetch_all(pool.get_ref()).await;

    match result {
        Ok(countries) => {
            let mut country_names: Vec<String> = countries.into_iter().filter_map(|c| c.country).collect();
            country_names.insert(0, "all".to_string());
            HttpResponse::Ok().json(json!({"countries": country_names}))
        },
        Err(e) => {
            eprintln!("Failed to get countries: {}", e);
            HttpResponse::InternalServerError().json(json!({"error": "Failed to retrieve countries"}))
        }
    }
}


#[derive(Deserialize)]
struct ForeignAidRequest {
    country: String,
    year: String
}
#[post("/foreign_aid")]
async fn post_foreign_aid(filters: web::Json<ForeignAidRequest>, pool: web::Data<MySqlPool>) -> impl Responder {
    let ForeignAidRequest {country, year} = filters.into_inner();

    println!("{}, {}", country, year);
    let results = make_map_data(&year, &pool).await;
    HttpResponse::Ok().json(results)
}

#[get("/comparison")]
async fn get_comparison() -> impl Responder {
    println!("GET comparison");
    HttpResponse::Ok().json(json!({"data": "some data"}))
}

#[get("/info")]
async fn get_info() -> impl Responder {
    println!("GET info");
    HttpResponse::Ok().json(json!({"data": "some data"}))
}
