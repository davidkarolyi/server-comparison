use actix_web::{web, App, HttpRequest, HttpServer, Responder};
use jsonwebtoken::{decode, errors, DecodingKey, Validation};
use serde::{Deserialize, Serialize};
use serde_json;
use std::fs;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
  println!("rust_actix is listening on localhost:3000");
  HttpServer::new(|| App::new().service(web::resource("/").to(handler)))
    .bind("0.0.0.0:3000")?
    .run()
    .await
}

async fn handler(_req: HttpRequest) -> impl Responder {
  let test_data = read_test_data().await;
  let payload = verify_token(test_data.token, test_data.secret).expect("Cannot decode token");
  web::Json(payload)
}

async fn read_test_data() -> TestData {
  let json_test_data = fs::read_to_string("../test_data.json").expect("Cannot read file");
  let test_data: TestData = serde_json::from_str(&json_test_data).unwrap();
  test_data
}

fn verify_token(token: String, secret: String) -> Result<Payload, errors::Error> {
  let validation = &mut Validation::default();
  validation.validate_exp = false;

  let token = decode::<Payload>(
    &token,
    &DecodingKey::from_secret(secret.as_ref()),
    validation,
  )?;
  Ok(token.claims)
}

#[derive(Debug, Serialize, Deserialize)]
struct TestData {
  token: String,
  secret: String,
}

#[allow(non_snake_case)]
#[derive(Debug, Serialize, Deserialize)]
struct Payload {
  userID: String,
  role: String,
  age: i32,
  iat: i32,
}
