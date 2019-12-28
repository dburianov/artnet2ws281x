package main

import (
        "bytes"

        ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"

        "log"
        "net"
        "github.com/spf13/viper"
)

var (
        brightness = 63
        width      = 150
        ledCounts  = 1200
        ArtNetHeaderSize = 18
        GitCommit        string
)

var ArtNetCheck, isOpcodeDmx bool = false, false
var (
        SubnetID1   = 0
        SubnetID2   = 0
        SubnetID3   = 0
        SubnetID4   = 0
        SubnetID5   = 0
        SubnetID6   = 0
        SubnetID7   = 0
        SubnetID8   = 0
        UniverseID1 = 1
        UniverseID2 = 2
        UniverseID3 = 3
        UniverseID4 = 4
        UniverseID5 = 5
        UniverseID6 = 6
        UniverseID7 = 7
        UniverseID8 = 8
)

func init() {
        viper.SetConfigType("yaml")                 // or viper.SetConfigType("YAML")
        viper.SetConfigName("config")               // name of config file (without extension)
        viper.AddConfigPath("/etc/artnet2ws281x/")  // path to look for the config file in
        viper.AddConfigPath("$HOME/.artnet2ws281x") // call multiple times to add many search paths
        viper.AddConfigPath(".")    // optionally look for config in the working directory
        err := viper.ReadInConfig() // Find and read the config file
        if err != nil {             // Handle errors reading the config file
                //              panic(fmt.Errorf("Fatal error config file: %s \n", err))
                log.Println("No config file: ", err)

                var yamlConfig = []byte(`
                        brightness = 63
                        SubnetID1: 0
                        SubnetID2: 0
                        SubnetID3: 0
                        SubnetID4: 0
                        SubnetID5: 0
                        SubnetID6: 0
                        SubnetID7: 0
                        SubnetID8: 0
                        UniverseID1: 1
                        UniverseID2: 2
                        UniverseID3: 3
                        UniverseID4: 4
                        UniverseID5: 5
                        UniverseID6: 6
                        UniverseID7: 7
                        UniverseID8: 8
                        `)

                viper.ReadConfig(bytes.NewBuffer(yamlConfig))
                //              return *debugmode
        }
        brightness = viper.GetInt("brightness")
        SubnetID1 = viper.GetInt("SubnetID1")
        SubnetID2 = viper.GetInt("SubnetID2")
        SubnetID3 = viper.GetInt("SubnetID3")
        SubnetID4 = viper.GetInt("SubnetID4")
        SubnetID5 = viper.GetInt("SubnetID5")
        SubnetID6 = viper.GetInt("SubnetID6")
        SubnetID7 = viper.GetInt("SubnetID7")
        SubnetID8 = viper.GetInt("SubnetID8")
        UniverseID1 = viper.GetInt("UniverseID1")
        UniverseID2 = viper.GetInt("UniverseID2")
        UniverseID3 = viper.GetInt("UniverseID3")
        UniverseID4 = viper.GetInt("UniverseID4")
        UniverseID5 = viper.GetInt("UniverseID5")
        UniverseID6 = viper.GetInt("UniverseID6")
        UniverseID7 = viper.GetInt("UniverseID7")
        UniverseID8 = viper.GetInt("UniverseID8")

}

func checkError(err error) {
        if err != nil {
                panic(err)
        }
}

func main() {
        log.Println("!!! START artnet2ws281x !!! ")
        log.Println("Version: ", GitCommit)
        // git pull && export GIT_COMMIT=$(git rev-list -1 HEAD) && export DATE_BUILD=$(date +%Y:%m:%d:%H:%M:%S:%Z)&& export GitCommit=$DATE_BUILD"_:_"$GIT_COMMIT && go build -o artnet2ws281x  -ldflags "-X main.GitCommit=$GitCommit"

        ArtNetHead := []byte{65, 114, 116, 45, 78, 101, 116, 0} //Art-Net
        selectUniverse1 := ((SubnetID1 * 16) + UniverseID1)
        selectUniverse2 := ((SubnetID2 * 16) + UniverseID2)
        selectUniverse3 := ((SubnetID3 * 16) + UniverseID3)
        selectUniverse4 := ((SubnetID4 * 16) + UniverseID4)
        selectUniverse5 := ((SubnetID5 * 16) + UniverseID5)
        selectUniverse6 := ((SubnetID6 * 16) + UniverseID6)
        selectUniverse7 := ((SubnetID7 * 16) + UniverseID7)
        selectUniverse8 := ((SubnetID8 * 16) + UniverseID8)
        opt := ws2811.DefaultOptions
        opt.Channels[0].Brightness = brightness
        opt.Channels[0].LedCount = ledCounts

        dev, err := ws2811.MakeWS2811(&opt)
        checkError(err)

        checkError(dev.Init())
        defer dev.Fini()

        // listen to incoming udp packets
        pc, err := net.ListenPacket("udp", ":6454")
        if err != nil {
                log.Fatal(err)
        }
        defer pc.Close()

        for {
                buf := make([]byte, 550)
                pc.ReadFrom(buf)
                if err != nil {
                        continue
                }

                ArtNet2Leds := func(currentUniverse int) {
                        for i := 0; i < width*3; i = i + 3 {
                                color := uint32((int(buf[i+ArtNetHeaderSize]) * 256 * 256) + (int(buf[i+1+ArtNetHeaderSize]) * 256) + (int(buf[i+2+ArtNetHeaderSize])))
                                dev.Leds(0)[150*(currentUniverse-1)+i/3] = color
                        }
                }

                for i := 0; i < 8; i++ {
                        if buf[i] == ArtNetHead[i] {
                                ArtNetCheck = true
                        }
                }
                if ArtNetCheck {
                        if buf[9] == 80 && buf[8] == 0 {
                                isOpcodeDmx = true
                        } else {
                                isOpcodeDmx = false
                        }
                }
                if ArtNetCheck && isOpcodeDmx {
                        if selectUniverse1 == int(buf[14]) {
                                go ArtNet2Leds(1)
                        }
                        if selectUniverse2 == int(buf[14]) {
                                go ArtNet2Leds(2)
                        }
                        if selectUniverse3 == int(buf[14]) {
                                go ArtNet2Leds(3)
                        }
                        if selectUniverse4 == int(buf[14]) {
                                go ArtNet2Leds(4)
                        }
                        if selectUniverse5 == int(buf[14]) {
                                go ArtNet2Leds(5)
                        }
                        if selectUniverse6 == int(buf[14]) {
                                go ArtNet2Leds(6)
                        }
                        if selectUniverse7 == int(buf[14]) {
                                go ArtNet2Leds(7)
                        }
                        if selectUniverse8 == int(buf[14]) {
                                go ArtNet2Leds(8)
                        }
                }
                checkError(dev.Render())
        }
}
