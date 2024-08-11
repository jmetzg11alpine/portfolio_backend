use std::collections::HashMap;

pub fn get_renaming_map() -> HashMap<&'static str, &'static str> {
    [
        ("Department of the Treasury", "Treasury"),
        (
            "Department of Health and Human Services",
            "Health and Human",
        ),
        ("Department of Defense", "Defense"),
        ("Social Security Administration", "Social Security"),
        ("Department of Veterans Affairs", "Veterans Affairs"),
        ("Department of Agriculture", "Agriculture"),
        ("Office of Personnel Management", "OPM"),
        ("Department of Housing and Urban Development", "Housing"),
        ("Department of Transportation", "Transportation"),
        ("Department of Homeland Security", "Homeland Security"),
        ("Department of Energy", "Energy"),
        ("Department of Commerce", "Commerce"),
        ("Department of Education", "Education"),
        ("Environmental Protection Agency", "Environmental"),
        ("Department of the Interior", "Interior"),
        ("Department of State", "State"),
        ("General Services Administration", "General Services"),
        ("Department of Justice", "Justice"),
        ("Department of Labor", "Labor"),
        ("Pension Benefit Guaranty Corporation", "Pension"),
    ]
    .iter()
    .cloned()
    .collect()
}

pub fn get_colors() -> Vec<&'static str> {
    vec![
        "0, 128, 128",  // Teal
        "255, 99, 71",  // Tomato
        "124, 252, 0",  // Lawn Green
        "70, 130, 180", // Steel Blue
        "255, 215, 0",  // Gold
        "0, 191, 255",  // Deep Sky Blue
        "255, 69, 0",   // Orange Red
        "138, 43, 226", // Blue Violet
        "60, 179, 113", // Medium Sea Green
        "218, 165, 32", // Golden Rod
    ]
}
