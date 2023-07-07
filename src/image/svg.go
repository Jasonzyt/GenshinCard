package image

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/jasonzyt/genshincard/net"
)

const SVG_CONTENT_FORMAT = `
<svg width="360" height="210" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <style>

        @font-face {
            font-family: 'HYWenHei-85W';
            src: url('https://raw.githubusercontent.com/Jasonzyt/genshin-stats/main/assets/font/subset-HYWenHei-HEW.woff2') format('woff2'), url('https://raw.githubusercontent.com/Jasonzyt/genshin-stats/main/assets/font/subset-HYWenHei-HEW.woff') format('woff');
            font-weight: 900;
            font-style: normal;
            font-display: swap;
        }

        @keyframes fadein {
            to {
                opacity: 1;
            }
        }

        svg {
            box-shadow: 5px 8px 4px 0 #bbb;
            border-radius: 20px;
        }

        text,
        td,
        th {
            opacity: 0;
            animation: fadein 0.7s ease-in-out;
            animation-fill-mode: forwards;
        }

        * {
            font-family: 'HYWenHei-85W', sans-serif;
            font-weight: 900;
        }

        th {
            font-size: 16px;
        }

        td {
            font-size: 24px;
        }

        tr,
        td,
        th {
            padding: 0;
        }

    </style>

    <image id="background" xlink:href="%s" width="100%%" height="100%%" preserveAspectRatio="xMaxYMid slice"/>
    <text x="30px" y="45px" width="300px" fill="white" style="animation-delay: 200ms">
        <tspan font-size="24px">%s
            <tspan font-size="16px" dx="15px">Lv.<tspan font-size="20px" dx="2px">%d</tspan>
            </tspan>
        </tspan>
    </text>
    <text x="30px" y="68px" width="300px" fill="white" style="animation-delay: 400ms">
        <tspan font-size="16px">UID: %s</tspan>
    </text>
    <foreignObject x="20px" y="76px" width="320px" height="100px">
        <table xmlns="http://www.w3.org/1999/xhtml" style="width: 320px; color: white; line-height: 1.5em; text-align: center;">
            <tr>
                <th style="animation-delay: 500ms">活跃天数</th>
                <th style="animation-delay: 700ms">深境螺旋</th>
                <th style="animation-delay: 900ms">获得角色</th>
                <th style="animation-delay: 1100ms">达成成就</th>
            </tr>
            <tr>
                <td style="animation-delay: 600ms">%d</td>
                <td style="animation-delay: 800ms">%s</td>
                <td style="animation-delay: 1000ms">%d</td>
                <td style="animation-delay: 1200ms">%d</td>
            </tr>
            <tr>
                <th style="animation-delay: 550ms">开启宝箱</th>
                <th style="animation-delay: 750ms">供奉神瞳</th>
                <th style="animation-delay: 950ms">激活锚点</th>
                <th style="animation-delay: 1150ms">洞天仙力</th>
            </tr>
            <tr>
                <td style="animation-delay: 650ms">%d</td>
                <td style="animation-delay: 850ms">%d</td>
                <td style="animation-delay: 1050ms">%d</td>
                <td style="animation-delay: 1250ms">%d</td>
            </tr>
        </table>
    </foreignObject>

</svg>
`

func randomBackground() string {
	backgroundPath := "./assets/img/"
	backgroundDirEntry, err := os.ReadDir(backgroundPath)
	if err != nil {
		panic(err)
	}
	backgrounds := make([]string, 0, len(backgroundDirEntry))
	for _, file := range backgroundDirEntry {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			backgrounds = append(backgrounds, file.Name())
		}
	}
	rand.Seed(time.Now().UnixNano())
	return "https://raw.githubusercontent.com/Jasonzyt/genshin-stats/main/assets/img/" + backgrounds[rand.Intn(len(backgrounds))]
}

func GenerateSvg(profile *net.PlayerProfile, writer io.Writer) {
	stats := profile.Statistics
	chestNumber := stats.CommonChestNumber + stats.ExquisiteChestNumber + stats.PreciousChestNumber + stats.LuxuriousChestNumber + stats.MagicChestNumber
	culusNumbewr := stats.AnemoculusNumber + stats.GeoculusNumber + stats.ElectroculusNumber + stats.DendroculusNumber
	maxComfortNumber := 0
	for _, home := range profile.Homes {
		if home.ComfortNumber > maxComfortNumber {
			maxComfortNumber = home.ComfortNumber
		}
	}

	out := fmt.Sprintf(SVG_CONTENT_FORMAT,
		randomBackground(),
		profile.Role.NickName, profile.Role.Level, profile.RoleId,
		stats.ActiveDayNumber,
		stats.SpiralAbyss,
		stats.AvatarNumber,
		stats.AchievementNumber,

		chestNumber,
		culusNumbewr,
		stats.WayPointNumber,
		maxComfortNumber,
	)
	writer.Write([]byte(out))
}
