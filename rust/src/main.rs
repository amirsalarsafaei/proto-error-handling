
use std::net::SocketAddr;
use tonic::transport::Server;

mod userservice;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Define the address to host the gRPC server
    let addr: SocketAddr = "[::1]:8000".parse()?;
    
    let mut server = Server::builder();
    
    println!("gRPC server starting on {}", addr);
    
    let user_servicer = userservice::UserServiceImpl::new(userservice::UserRepository::new());

    server
        .add_service(user_servicer.into_service())
        .serve(addr)
        .await?;

    Ok(())
}

