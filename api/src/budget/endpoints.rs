use actix_web::{post, get, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use sqlx::{MySqlPool, query_as};
use serde_json::json;
use crate::budget::helpers::{process_agency_data, AgencyBudget, make_map_data, make_bar_data, MapResults, BarResults};


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

#[derive(Deserialize)]
struct ForeignAidRequest {
    country: String,
    year: String
}

#[derive(Serialize)]
struct ForeignAidResponse {
    map_results: MapResults,
    bar_results: BarResults
}
#[post("/foreign_aid")]
async fn post_foreign_aid(filters: web::Json<ForeignAidRequest>, pool: web::Data<MySqlPool>) -> impl Responder {
    let ForeignAidRequest {country, year} = filters.into_inner();

    let map_results = make_map_data(&year, &pool).await;
    let bar_results = make_bar_data(&year, &country, &pool).await;

    let foreign_aid_response = ForeignAidResponse {
        map_results,
        bar_results
    };
    HttpResponse::Ok().json(foreign_aid_response)
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
