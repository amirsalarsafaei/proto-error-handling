use proto_error_interface::hello_world::{
    create_user_alt_response::Result as AltResult,
    user_service_server::{UserService, UserServiceServer},
    CreateUserAltResponse, CreateUserRequest, CreateUserResponse,
    ErrorDetails as InternalErrorDetails, UserData, UserStatus,
};
use std::{collections::HashMap, sync::Arc};
use tokio::sync::Mutex;
use tonic::{Code, Request, Response, Status};

use thiserror::Error;
use tonic_types::{ErrorDetails, StatusExt};
use uuid::Uuid;

#[derive(Error, Debug)]
pub enum UserError {
    #[error("Invalid email format: {0}")]
    InvalidEmail(String),
    #[error("Invalid username: {0}")]
    InvalidUsername(String),
}

#[derive(Clone)]
pub struct User {
    id: Uuid,
    email: String,
    username: String,
}

impl User {
    pub fn new(email: String, username: String) -> Result<Self, UserError> {
        if !Self::is_valid_email(&email) {
            return Err(UserError::InvalidEmail(email));
        }

        if !Self::is_valid_username(&username) {
            return Err(UserError::InvalidUsername(
                "Username must be between 3 and 30 characters and contain only alphanumeric characters and underscores"
                    .to_string(),
            ));
        }

        Ok(Self {
            id: Uuid::new_v4(),
            email,
            username,
        })
    }

    fn is_valid_email(email: &str) -> bool {
        email_address::EmailAddress::is_valid(email)
    }

    fn is_valid_username(username: &str) -> bool {
        // Username requirements:
        // 1. 3-30 characters
        // 2. Only alphanumeric and underscore
        // 3. Can't start with number
        let username_regex = regex::Regex::new(r"^[a-zA-Z][a-zA-Z0-9_]{2,29}$").unwrap();
        username_regex.is_match(username)
    }
}

#[derive(Clone)]
pub struct UserRepository {
    users: Arc<Mutex<Vec<User>>>,
}

impl UserRepository {
    pub fn new() -> Self {
        UserRepository {
            users: Arc::new(Mutex::new(Vec::new())),
        }
    }

    pub async fn create(&self, user: User) -> Result<User, Status> {
        let mut users = self.users.lock().await;

        users.push(user.clone());
        Ok(user)
    }

    pub async fn exists(&self, user: &User) -> bool {
        self.users
            .lock()
            .await
            .iter()
            .any(|u| u.id == user.id || u.email == user.email || u.username == user.username)
    }
}

pub struct UserServiceImpl {
    repository: UserRepository,
}

impl UserServiceImpl {
    pub fn new(repository: UserRepository) -> Self {
        UserServiceImpl { repository }
    }

    pub fn into_service(self) -> UserServiceServer<Self> {
        UserServiceServer::new(self)
    }
}

#[tonic::async_trait]
impl UserService for UserServiceImpl {
    async fn create_user(
        &self,
        request: Request<CreateUserRequest>,
    ) -> Result<Response<CreateUserResponse>, Status> {
        let user_request = request.into_inner();
        let user = User::new(user_request.email, user_request.username).map_err(|e| match e {
            UserError::InvalidEmail(msg) => Status::with_error_details(
                Code::InvalidArgument,
                "Validation error",
                ErrorDetails::with_bad_request_violation("email", msg),
            ),
            UserError::InvalidUsername(msg) => Status::with_error_details(
                Code::InvalidArgument,
                "Validation error",
                ErrorDetails::with_bad_request_violation("username", msg),
            ),
        })?;

        if self.repository.exists(&user).await {
            return Err(Status::with_error_details(
                Code::AlreadyExists,
                "Resource already exists",
                ErrorDetails::with_bad_request_violation(
                    "user",
                    "User with this email or username already exists",
                ),
            ));
        }

        match self.repository.create(user).await {
            Ok(user) => Ok(Response::new(CreateUserResponse {
                user_id: user.id.to_string(),
                status: UserStatus::Pending as i32,
            })),
            Err(e) => Err(Status::with_error_details(
                Code::Internal,
                "could not create user",
                ErrorDetails::with_error_info(
                    "could not create user",
                    "hello_world.UserService",
                    vec![
                        ("error_type".to_string(), "database_error".to_string()),
                        ("error_detail".to_string(), e.to_string()),
                    ]
                    .into_iter()
                    .collect::<HashMap<String, String>>(),
                ),
            )),
        }
    }

    async fn create_user_alt(
        &self,
        request: Request<CreateUserRequest>,
    ) -> Result<Response<CreateUserAltResponse>, Status> {
        let user_request = request.into_inner();

        let user = match User::new(user_request.email, user_request.username) {
            Ok(user) => user,
            Err(e) => {
                return Ok(Response::new(CreateUserAltResponse {
                    result: Some(AltResult::Error(InternalErrorDetails {
                        code: "VALIDATION_ERROR".to_string(),
                        message: e.to_string(),
                    })),
                }));
            }
        };

        // Check if user already exists
        if self.repository.exists(&user).await {
            return Ok(Response::new(CreateUserAltResponse {
                result: Some(AltResult::Error(InternalErrorDetails {
                    code: "ALREADY_EXISTS".to_string(),
                    message: "User with this email or username already exists".to_string(),
                })),
            }));
        }

        match self.repository.create(user).await {
            Ok(user) => Ok(Response::new(CreateUserAltResponse {
                result: Some(AltResult::Success(UserData {
                    user_id: user.id.to_string(),
                    status: UserStatus::Pending as i32,
                })),
            })),
            Err(e) => Ok(Response::new(CreateUserAltResponse {
                result: Some(AltResult::Error(InternalErrorDetails {
                    code: "CREATION_ERROR".to_string(),
                    message: format!("Failed to create user: {}", e),
                })),
            })),
        }
    }
}
