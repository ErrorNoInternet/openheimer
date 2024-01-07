// const WINPCAP_LIB_X32: &str = "libs/npcap/Lib/";
// const WINPCAP_LIB_ARM: &str = "libs/npcap/Lib/ARM64/";
const WINPCAP_LIB_X64: &str = "lib/npcap/Lib/x64/";

fn main() {
    println!("cargo:rerun-if-changed=build.rs");
    println!("cargo:rerun-if-changed=Cargo.lock");

    built::write_built_file().unwrap();

    let target = std::env::var("TARGET").unwrap();
    if target.contains("windows") && target.contains("x86_64") {
        println!("cargo:rustc-link-search=native={WINPCAP_LIB_X64}");
    }
}
