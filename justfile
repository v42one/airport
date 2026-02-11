k0sctl := "go tool k0sctl"

fmt:
    go tool fmt .

test:
    go test --count=1 ./...

apply cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} apply --config cluster.yaml

reset cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} reset --config cluster.yaml

kubeconfig cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} kubeconfig --config cluster.yaml > ~/.kube_config/config--{{ cluster }}.yaml

SINGBOX_VERSION := "v1.13.18"

dep-singbox:
    rm -rf target/sing-box
    git clone --recurse-submodules --shallow-submodules --depth=1 -b {{ SINGBOX_VERSION }} \
          git@github.com:SagerNet/sing-box.git ./target/sing-box

[working-directory('target/sing-box')]
build-sfm-lib:
    make lib_install
    make lib_apple

[working-directory('target/sing-box/clients/apple')]
build-sfm:
    cp -rf ../../Libbox.xcframework ./Libbox.xcframework
    rm -rf build/SFM.System.xcarchive
    xcodebuild \
        -scheme SFM.System \
        -configuration Debug \
        -archivePath build/SFM.System.xcarchive \
        ARCHS=arm64 \
        CODE_SIGN_IDENTITY="" \
        DEVELOPMENT_TEAM="" \
        CODE_SIGNING_REQUIRED=NO \
        CODE_SIGNING_ALLOWED=NO \
        archive