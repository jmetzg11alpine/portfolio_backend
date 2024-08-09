use actix_web::{post, get, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};



#[get("/budget")]
async fn get_budget() -> impl Responder {
    HttpResponse::Ok().json("return budget data")
}

#[derive(Deserialize, Serialize)]
struct BudgetRequest {
    key: String
}

#[post("/foreignAid")]
async fn post_foreign_aid(budget_request: web::Json<BudgetRequest>) -> impl Responder {
    HttpResponse::Ok().json(budget_request.into_inner())
}

#[get("/comparison")]
async fn get_comparison() -> impl Responder {
    HttpResponse::Ok().json("get foreign aid data")
}

#[get("/info")]
async fn get_info() -> impl Responder {
    HttpResponse::Ok().json("get info data")
}
