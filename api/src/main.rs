use actix_web::{get, App, HttpResponse, HttpServer, Responder};
use actix_cors::Cors;

mod budget;

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Server is Running")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .wrap(
                Cors::default()
                    .allowed_origin("http://localhost:5173")
                    .allowed_methods(vec!["GET", "POST"])
                    .allowed_headers(vec![actix_web::http::header::CONTENT_TYPE])
                    .max_age(3600)
            )
            .service(hello)
            .service(budget::endpoints::get_budget)
            .service(budget::endpoints::post_foreign_aid)
            .service(budget::endpoints::get_comparison)
            .service(budget::endpoints::get_info)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
