use actix_web::{get, App, HttpResponse, HttpServer, Responder, web};
use actix_cors::Cors;
use sqlx::mysql::MySqlPool;
use std::env;

mod budget;

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Server is Running")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let database_url = env::var("MYSQL_URL").expect("MYSQL_URL must be set");
    let pool = MySqlPool::connect(&database_url).await.expect("Failed to creat pool");

    HttpServer::new(move || {
        App::new()
            .wrap(
                Cors::default()
                    .allowed_origin("http://localhost:5173")
                    .allowed_methods(vec!["GET", "POST"])
                    .allowed_headers(vec![actix_web::http::header::CONTENT_TYPE])
                    .max_age(3600)
            )
            .app_data(web::Data::new(pool.clone()))
            .service(hello)
            .service(budget::endpoints::get_agency)
            .service(budget::endpoints::post_foreign_aid)
            .service(budget::endpoints::get_comparison)
            .service(budget::endpoints::get_info)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
