fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_prost_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(&["./proto/trajectory.proto"], &["./proto"])?;

    let build_tle_server = std::env::var("CARGO_FEATURE_GRPC_TEST_SERVER").is_ok();

    tonic_prost_build::configure()
        .build_server(build_tle_server)
        .build_client(true)
        .compile_protos(&["./proto/tle.proto"], &["./proto"])?;

    Ok(())
}
