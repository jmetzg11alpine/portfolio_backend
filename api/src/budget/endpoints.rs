use actix_web::{post, get, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use sqlx::MySqlPool;
use serde_json::json;
use crate::budget::data::{get_renaming_map, get_colors};


#[derive(Serialize)]
struct AgencyEntry {
    label: String,
    value: f32,
    tooltip: String,
    background_color: String,
}


#[get("/agency")]
async fn get_agency(pool: web::Data<MySqlPool>) -> impl Responder {
    let renaming = get_renaming_map();
    let colors = get_colors();

    let result = sqlx::query!("SELECT id, agency, budget FROM agency_budget").fetch_all(pool.get_ref()).await;

    match result {
        Ok(rows) => {
            let mut main_data = Vec::new();
            let mut other_data = Vec::new();
            let mut main_other_value = 0.0;
            let mut other_other_value = 0.0;
            let mut other_other_labels = Vec::new();

            for (i, row) in rows.iter().enumerate() {
                let agency = row.agency.as_deref().unwrap_or("unknown agency");
                let budget = row.budget.unwrap_or(0.0);
                let label = renaming.get(agency).unwrap_or(&agency).to_string();
                let background_color = format!("rgba({}, 1)", colors[i % colors.len()]);

                if i < 9 {
                    main_data.push(AgencyEntry{
                        label,
                        value: budget,
                        tooltip: agency.to_string(),
                        background_color: background_color
                    });
                } else if i < 18 {
                    main_other_value += budget;
                    other_data.push(AgencyEntry {
                        label,
                        value: budget,
                        tooltip: agency.to_string(),
                        background_color: background_color
                    });
                } else {
                    main_other_value += budget;
                    other_other_value += budget;
                    other_other_labels.push((agency.to_string(), budget));
                }
            }

            main_data.push(AgencyEntry {
                label: ("other").to_string(),
                value: main_other_value,
                tooltip: ("break down in other graph").to_string(),
                background_color: format!("rbga({}, 1)", colors[9])
            });

            other_data.push(AgencyEntry {
                label: format!("{} others", other_other_labels.len()),
                value: other_other_value,
                tooltip: ("break down described below").to_string(),
                background_color: format!("rbga({}, 1)", colors[9])
            });
            HttpResponse::Ok().json(json!({
                "main_data": main_data,
                "other_data": other_data,
                "other_other_labels": other_other_labels
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
