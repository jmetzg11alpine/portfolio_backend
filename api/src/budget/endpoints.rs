use actix_web::{post, get, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use sqlx::MySqlPool;
use sqlx::query_as;
use serde_json::json;
use crate::budget::helpers::{process_agency_data, AgencyBudget};


#[get("/agency")]
async fn get_agency(pool: web::Data<MySqlPool>) -> impl Responder {

    let result = query_as!(AgencyBudget, "Select agency, budget FROM agency_budget order by budget desc").fetch_all(pool.get_ref()).await;

    match result {
      Ok(rows) => {
        let (main_data, other_data, table_labels) = process_agency_data(&rows);
            HttpResponse::Ok().json(json!({
                "main_data": main_data,
                "other_data": other_data,
                "table_labels": table_labels
            }))
    }
         Err(e) => {
            eprintln!("Failed to execute query: {}", e);
            HttpResponse::InternalServerError().json(json!({"error": "Failed to retrieve data"}))
        }
    }
}


#[derive(Deserialize, Serialize)]
struct ForeignAidRequest {
    key: String
}

#[post("/foreignAid")]
async fn post_foreign_aid(budget_request: web::Json<ForeignAidRequest>) -> impl Responder {
    println!("post foreignAid");
    HttpResponse::Ok().json(budget_request.into_inner())
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
