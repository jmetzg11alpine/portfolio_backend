use crate::budget::data::{get_colors, get_renaming_map};
use serde::Serialize;

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
