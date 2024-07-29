use actix_web::{get, App, HttpResponse, HttpServer, Responder};

mod budget;

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Server is Running")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(hello)
            .service(budget::endpoints::get_budget)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
