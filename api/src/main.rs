use actix_web::{get, App, HttpResponse, HttpServer, Responder, web};
use actix_cors::Cors;
use sqlx::mysql::MySqlPool;
use std::env;
use dotenv::dotenv;

mod budget;

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Server is Running")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();

   let database_url = match env::var("DATABASE_URL") {
        Ok(url) => url,
        Err(e) => {
            eprintln!("Error reading DATABASE_URL: {:?}", e);
            return Err(std::io::Error::new(std::io::ErrorKind::Other, "DATABASE_URL not set"));
        }
    };

    let pool = match MySqlPool::connect(&database_url).await {
        Ok(pool) => pool,
        Err(e) => {
            eprintln!("Failed to create pool: {:?}", e);
            return Err(std::io::Error::new(std::io::ErrorKind::Other, "Database connection failed"));
        }
    };

    HttpServer::new(move || {
        println!("Server setup in progress...");
        App::new()
            .wrap(
                Cors::default()
                    .allowed_origin("http://localhost:5173")
                    .allowed_origin("http://localhost")
                    .allowed_origin("http://192.168.1.241")
                    .allowed_methods(vec!["GET", "POST"])
                    .allowed_headers(vec![actix_web::http::header::CONTENT_TYPE])
                    .max_age(3600)
            )
            .app_data(web::Data::new(pool.clone()))
            .service(hello)
            .service(budget::endpoints::get_agency)
            .service(budget::endpoints::post_foreign_aid)
            .service(budget::endpoints::get_comparison)
    })
    .bind(("0.0.0.0", 8080))
    .map_err(|e| {
        eprintln!("Failed to bind server: {:?}", e);
        e
    })?
    .run()
    .await
    .map_err(|e| {
        eprintln!("Server failed to run: {:?}", e);
        e
    })
}
