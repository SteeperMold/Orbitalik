fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_prost_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(&["./proto/trajectory.proto"], &["./proto"])?;

    tonic_prost_build::configure()
        .build_server(false)
        .build_client(true)
        .compile_protos(&["./proto/tle.proto"], &["./proto"])?;

    Ok(())
}
