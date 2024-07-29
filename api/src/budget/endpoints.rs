use actix_web::{get, HttpResponse, Responder};

#[get("/budget")]
async fn get_budget() -> impl Responder {
    HttpResponse::Ok().body("Budget endpoint")
}
