use crate::budget::data::{get_colors, get_renaming_map};
use serde::Serialize;
use sqlx::{FromRow, MySqlPool, query_as};
use std::collections::{HashMap, HashSet};

#[derive(Debug)]
pub struct AgencyBudget {
    pub agency: Option<String>,
    pub budget: Option<f32>,
}

#[derive(Serialize)]
pub struct AgencyEntry {
    label: String,
    value: f32,
    tooltip: String,
    background_color: String,
}

pub fn process_agency_data(
    rows: &[AgencyBudget],
) -> (Vec<AgencyEntry>, Vec<AgencyEntry>, Vec<(String, f32)>) {
    let renaming = get_renaming_map();
    let colors = get_colors();

    let mut main_data = Vec::new();
    let mut other_data = Vec::new();
    let mut main_other_value = 0.0;
    let mut other_other_value = 0.0;
    let mut table_data = Vec::new();

    for (i, row) in rows.iter().enumerate() {
        let agency = row.agency.as_deref().unwrap_or("unknown agency");
        let budget = row.budget.unwrap_or(0.0);
        let label = renaming.get(agency).unwrap_or(&agency).to_string();

        if i < 9 {
            main_data.push(AgencyEntry {
                label,
                value: budget,
                tooltip: agency.to_string(),
                background_color: format!("rgba({}, 1)", colors[i]),
            })
        } else if i < 18 {
            main_other_value += budget;
            other_data.push(AgencyEntry {
                label,
                value: budget,
                tooltip: agency.to_string(),
                background_color: format!("rgba({}, 1)", colors[i - 9]),
            })
        } else {
            main_other_value += budget;
            other_other_value += budget;
            table_data.push((agency.to_string(), budget));
        }
    }
    main_data.push(AgencyEntry {
        label: ("other").to_string(),
        value: main_other_value,
        tooltip: ("break down in other graph").to_string(),
        background_color: format!("rgba({}, 1)", colors[9]),
    });

    other_data.push(AgencyEntry {
        label: format!("{} others", table_data.len()),
        value: other_other_value,
        tooltip: ("break down described below").to_string(),
        background_color: format!("rgba({}, 1)", colors[9]),
    });

    (main_data, other_data, table_data)
}

#[derive(FromRow)]
struct ForeignAidMapResults {
    country: String,
    amount: f32,
    lat: f32,
    lng: f32
}

#[derive(Serialize)]
struct MapData {
    lat: f32,
    lng: f32,
    text: String,
    size: f32,
    amount: f32,
}

#[derive(Serialize)]
pub struct MapResults {
    map_data: Vec<MapData>,
    countries: Vec<String>
}

pub async fn make_map_data(year: &str, pool: &MySqlPool) -> MapResults {
    let year_query = if year != "all" {
        format!("SELECT * FROM foreign_aid WHERE year = {}", year)
    } else {
        String::from("SELECT * FROM foreign_aid")
    };
   let results = sqlx::query_as::<_, ForeignAidMapResults>(&year_query).fetch_all(pool).await.expect("Failed to fetch map data");

    prep_map_data(results, &year)
}


fn prep_map_data(results: Vec<ForeignAidMapResults>, year: &str) -> MapResults {
    let mut country_map: HashMap<String, MapData> = HashMap::new();
    let mut unique_countries: HashSet<String> = HashSet::new();

    for row in results.iter() {
        let country = row.country.clone();
        unique_countries.insert(country.clone());

        if let Some(map_data) = country_map.get_mut(&country) {
            map_data.amount += row.amount;
        } else {
            let new_map_data = MapData {
                lat: row.lat,
                lng: row.lng,
                text: country.clone(),
                size: 0.0,
                amount: row.amount,
            };
            country_map.insert(country, new_map_data);
        }

    }

    let min_amount = country_map.values().map(|data| data.amount).fold(f32::MAX, f32::min);
    let max_amount = country_map.values().map(|data| data.amount).fold(f32::MIN, f32::max);

    for map_data in country_map.values_mut() {
        map_data.size = normalize_amount(map_data.amount, min_amount, max_amount);
        map_data.text = make_map_text(year, &map_data.text);
    }

    let map_data = country_map.into_values().collect();

    let mut countries: Vec<String> = unique_countries.into_iter().collect();
    countries.sort();
    countries.insert(0, "all".to_string());

    MapResults {map_data, countries}
}

fn normalize_amount(x: f32, min: f32, max: f32) -> f32 {
    let base_size = 4.0;
    let scale_factor = 35.0;
    if x <= min {
        return base_size;
    }
    base_size + scale_factor * (x - min) / (max - min)
}

fn make_map_text(year: &str, country: &str) -> String {
    if year == "all" {
        format!("{} (10 yrs.): ", country)
    } else {
        format!("{} in {}: ", country, year)
    }
}

#[derive(FromRow)]
struct ForeignAidBarQuery{
    amount: f32,
    year: i32
}

#[derive(Serialize)]
struct BarData {
    year: i32,
    amount: f32,
}

#[derive(Serialize)]
pub struct BarResults {
    bar_data: Vec<BarData>,
    total_amount: f32
}


pub async fn make_bar_data(year: &str, country: &str, pool: &MySqlPool) -> BarResults {
    let country_query = if country != "all" {
        format!("SELECT amount, year FROM foreign_aid WHERE country = '{}'", country)
    } else {
        String::from("SELECT amount, year from foreign_aid")
    };
    let results = sqlx::query_as::<_, ForeignAidBarQuery>(&country_query)
        .fetch_all(pool)
        .await
        .expect("Failed to fetch bar data");

    let mut bar_data_map: HashMap<i32, f32> = (2015..2025)
        .map(|year| (year as i32, 0.0))
        .collect();

    for aid in results {
        *bar_data_map.entry(aid.year).or_insert(0.0) += aid.amount;
    }

    let mut bar_data: Vec<BarData> = bar_data_map
        .into_iter()
        .map(|(year, amount)| BarData {year, amount})
        .collect();

    bar_data.sort_by_key(|data| data.year);

    let total_amount: f32 = if year == "all" {
        bar_data.iter()
            .map(|data| data.amount)
            .sum()
    } else {
        let specific_year: i32 = year.parse().expect("Failed to parse year");
        bar_data
            .iter()
            .filter(|data| data.year == specific_year)
            .map(|data| data.amount)
            .sum()
    };

    BarResults {
        bar_data,
        total_amount
    }
}


#[derive(Serialize)]
struct ComparisonQuery {
    year: Option<i32>,
    name: Option<String>,
    amount: Option<f32>,
}

#[derive(Serialize)]
pub struct ComparisonResults {
    data: HashMap<String, Vec<(i32, i32)>>,
    x_labels: Vec<i32>,
    agencies: Vec<String>,
}

pub async fn make_comparison_data(pool: &MySqlPool) -> ComparisonResults {
    let results = query_as!(
        ComparisonQuery,
        "SELECT year, name, amount FROM function_spending"
    )
    .fetch_all(pool)
    .await
    .expect("Failed to get comparison function spending data");

    let mut data: HashMap<String, Vec<(i32, i32)>> = HashMap::new();
    let mut x_labels: HashSet<i32> = HashSet::new();
    let mut agencies: HashSet<String> = HashSet::new();

    for result in results {
        if let(Some(year), Some(name), Some(amount)) = (result.year, result.name, result.amount) {
            let value = (amount / 1_000_000_000.0).round() as i32;
            data.entry(name.clone()).or_insert_with(Vec::new).push((year, value));
            x_labels.insert(year);
            agencies.insert(name);
        }
    }

    let mut x_labels: Vec<i32> = x_labels.into_iter().collect();
    x_labels.sort();

    let mut agencies: Vec<String> = agencies.into_iter().collect();
    agencies.sort();

    for values in data.values_mut() {
        values.sort_by(|a, b| a.0.cmp(&b.0));
    }

    ComparisonResults {
        data,
        x_labels,
        agencies
    }
}
