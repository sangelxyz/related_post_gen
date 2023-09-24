git checkout go-pgo &&
    ./gen_prof.sh &&
    cp go_gen_prof/cpu.prof /tmp/ &&
    git checkout main &&
    cp /tmp/cpu.prof go/
