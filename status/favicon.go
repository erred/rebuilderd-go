package status

import (
	"encoding/base64"
	"net/http"
)

func favicon(w http.ResponseWriter, r *http.Request) {
	b, err := base64.StdEncoding.DecodeString(favi)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "image/x-icon")
	w.Write(b)
}

const favi = `
AAABAAEAICAAAAEAIACoEAAAFgAAACgAAAAgAAAAQAAAAAEAIAAAAAAAABAAACMuAAAjLgAAAAAA
AAAAAAASEg7/S0tE/3Nvav9jYFr/Z2Rf/2dkXv9nZF7/Z2Re/2dkXv9nZF7/Z2Re/2dkXv9nZF7/
Z2Re/2dkXv9nZF7/Z2Re/2dkXv9mZF7/aWZg/15eV/86Ki3/ZzhL/2Q4Sf9kOEn/ZDhJ/2Q4Sf9k
OEn/ZDhJ/2Q3SP9nOUv/NiQo/1sxQf9pTFX/xcm8//Hr4//h3dT/5eHY/+Xh2P/l4dj/5eHY/+Xh
2P/l4dj/5eHY/+Xh2P/l4dj/5eHY/+Xh2P/l4dj/5eHY/+Tg1//r5d3/zNDD/3NHV//xcKv/526l
/+lvpv/pb6b/6W+m/+lvpv/pb6b/526l//FyrP9nOUv/fENa/9NflP9INjr/vL+z/+fg2P/W0sr/
29fO/9rWzf/a1s3/2tbN/9rWzf/a1s3/2tbN/9rWzf/a1s3/2tbN/9rWzf/a1s3/2dXM/+Da0v/D
xrr/bkVU/+Zro//cap7/3mqf/95qn//eap//3mqf/95qn//caZ3/526l/2Q3SP9tO0//9Xaw/8db
jP9KODz/v8G2/+ni2v/Y1Mv/3NjP/9zYz//c2M//3NjP/9zYz//c2M//3NjP/9zYz//c2M//3NjP
/9zYz//b187/4tzU/8TIvP9vRVT/6Gyl/95rn//ga6D/4Gug/+BroP/ga6D/4Gug/95qn//pb6b/
ZDhJ/3E9Uv/mbaT/6XGn/8hcjP9IODr/v8G2/+ji2v/Y1Mv/3NjP/9zYz//c2M//3NjP/9zYz//c
2M//3NjP/9zYz//c2M//3NjP/9vXzv/i3NT/xMi8/29FVP/obKX/3muf/+BroP/ga6D/4Gug/+Br
oP/ga6D/3mqf/+lvpv9kOEn/cT1S/+pvp//aaJz/7HKp/8Raiv9GNzn/wMK3/+ji2v/Y1Mv/3NjP
/9zYz//c2M//3NjP/9zYz//c2M//3NjP/9zYz//c2M//29fO/+Lc1P/EyLz/b0VU/+hspf/ea5//
4Gug/+BroP/ga6D/4Gug/+BroP/eap//6W+m/2Q4Sf9xPVL/6m+n/95qn//caZ3/7XKq/8NZif9C
NTb/wcK3/+ji2v/Y1Mv/3NjP/9zYz//c2M//3NjP/9zYz//c2M//3NjP/9zYz//b187/4tzU/8TI
vP9vRVT/6Gyl/95rn//ga6D/4Gug/+BroP/ga6D/4Gug/95qn//pb6b/ZDhJ/3E9Uv/qb6f/3mqe
/+BroP/caZ3/7XKq/8BZh/9BNTX/wcK3/+ji2v/Y1Mv/3NjP/9zYz//c2M//3NjP/9zYz//c2M//
3NjP/9vXzv/i3NT/xMi8/29FVP/obKX/3muf/+BroP/ga6D/4Gug/+BroP/ga6D/3mqf/+lvpv9k
OEn/cT1S/+pvp//eap7/4Gug/+BroP/caZ3/7nOr/71Xhf8/NDT/wsK4/+ji2v/Y1Mv/3NjP/9zY
z//c2M//3NjP/9zYz//c2M//29fO/+Lc1P/FyLz/bkVU/+Vro//cap3/3Wqe/91qnv/dap7/3Wqe
/91qnv/caZ3/5m2k/2M3SP9xPVL/6m+n/95qnv/ga6D/4Gug/+BroP/caZ3/7nOr/7tWhP88MzL/
wsO4/+jj2v/Y1Mv/3NjP/9zYz//c2M//3NjP/9zYz//b187/4tzU/8TIu/9zR1f/9XKu/+pwqP/s
cKn/7HCp/+xwqf/scKn/7HCp/+pvp//1dK//aTpM/3E9Uv/qb6f/3mqe/+BroP/ga6D/4Gug/+Br
oP/caZ3/73Or/7hVgf88MTH/wb+2/+jj2v/Y1Mv/3NjP/9zYz//c2M//3NjP/9vXzv/g29L/zcvB
/0IyNf9gM0X/XzZG/182Rf9fNkX/XzZF/182Rf9fNkX/XzVF/2I3R/80JCf/cT1S/+pvp//eap7/
4Gug/+BroP/ga6D/4Gug/+BroP/baZ3/8nKs/2o+T/8+RDv/wr21/+jk2//Y1Mv/3NjP/9zYz//c
2M//29fO/9/b0v/SzsX/Jich/woPCP8PFQ3/DhQM/w4UDP8OFAz/DhQM/w4UDP8OFAz/DhMM/xYX
Ev9xPVL/6m+n/95qnv/ga6D/4Gug/+BroP/ga6D/4Gug/99roP/kaKH/fFNi/7K2qv9RTUj/vrqy
/+nk2//Y1Mv/3NjP/9zYz//b187/39vS/9HNxP8sKSX/GhYT/x8bGP8eGhf/HhoX/x4aF/8eGhf/
HhoX/x4aF/8eGhf/HBkW/3E9Uv/qb6f/3mqe/+BroP/ga6D/4Gug/+BroP/ga6D/3muf/+drpP9x
R1b/1NfL/8vGvv9LSUT/v7uz/+nk2//Y1Mv/3NjP/9vXzv/f29L/0c3F/ysoJP8XFRH/HBoW/xsZ
Ff8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/cT1S/+pvp//eap7/4Gug/+BroP/ga6D/4Gug
/+BroP/fa5//52uk/3RLWv/Fyb3/6+bd/8bCuv9MSkT/v7uz/+nk2//Y1Mv/3NjP/9/b0v/RzcX/
Kygk/xcVEf8cGhb/GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf9xPVL/6m+n/95qnv/g
a6D/4Gug/+BroP/ga6D/4Gug/99rn//na6T/dEpZ/8nNwP/d18//5uLZ/8bDuv9MSkT/v7uz/+nk
2//X08v/39vS/9HNxf8rKCT/FxUR/xwaFv8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/GxkV
/3E9Uv/qb6f/3mqe/+BroP/ga6D/4Gug/+BroP/ga6D/32uf/+drpP90Sln/yczA/+Hb0//Y1Mv/
5+PZ/8bDuv9MSkT/v7uz/+jk2v/b187/0s7F/ysoJP8XFRH/HBoW/xsZFf8bGRX/GxkV/xsZFf8b
GRX/GxkV/xsZFf8bGRX/cT1S/+pvp//eap7/4Gug/+BroP/ga6D/4Gug/+BroP/fa5//52uk/3RK
Wf/JzMD/4dvT/9zYz//Y1Mz/5+PZ/8bDuv9MSkT/vrqy/+vn3v/NycH/Kykl/xcVEf8cGhb/GxkV
/xsZFf8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf9xPVL/6m+n/95qnv/ga6D/4Gug/+BroP/ga6D/
4Gug/99rn//na6T/dEpZ/8nMwP/h29P/29fO/9zYz//Y1Mz/5+PZ/8bDuv9LSUT/wb61/97a0f8o
JiH/FxUS/xwaFv8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/GxkV/3E9Uv/qb6f/3mqe/+Br
oP/ga6D/4Gug/+BroP/ga6D/32uf/+drpP90Sln/yczA/+Hb0//b187/3NjP/9zYz//Y1Mz/5+PZ
/8bCuf9PTUj/sq6n/zQyLf8VEw//HBoW/xsZFf8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/
cT1S/+pvp//eap7/4Gug/+BroP/ga6D/4Gug/+BroP/fa5//52uk/3RKWf/JzMD/4dvT/9vXzv/c
2M//3NjP/9zYz//Y1Mz/5uLZ/8fDu/9KSEP/HBoW/xwaFv8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZ
Ff8bGRX/GxkV/xsZFf9xPVL/6m+n/95qnv/ga6D/4Gug/+BroP/ga6D/4Gug/99rn//na6T/dEpZ
/8nMwP/h29P/29fO/9zYz//c2M//3NjP/9zYz//Y1cz/5eHY/8zIwP81My7/EQ8L/x8dGf8bGRX/
GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/GxkV/3E9Uv/qb6f/3mqe/+BroP/ga6D/4Gug/+BroP/g
a6D/32uf/+drpP90Sln/yczA/+Hb0//b187/3NjP/9zYz//c2M//3NjP/9zYz//Z1cz/5ODX/9DM
w/83NTD/EA4L/x8dGf8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf8bGRX/cT1S/+pvp//eap7/4Gug
/+BroP/ga6D/4Gug/+BroP/fa5//52uk/3RKWf/JzMD/4dvT/9vXzv/c2M//3NjP/9zYz//c2M//
3NjP/9zYz//Z1cz/49/W/9DMw/83NTD/EA4L/x8dGf8bGRX/GxkV/xsZFf8bGRX/GxkV/xsZFf9x
PVL/6m+n/95qnv/ga6D/4Gug/+BroP/ga6D/4Gug/99rn//na6T/dEpZ/8nMwP/h29P/29fO/9zY
z//c2M//3NjP/9zYz//c2M//3NjP/9zYz//Z1cz/49/W/9DMw/83NTD/EA4L/x8dGf8bGRX/GxkV
/xsZFf8bGRX/GxkV/3E9Uv/qb6f/3mqe/+BroP/ga6D/4Gug/+BroP/ga6D/32uf/+drpP90Sln/
yczA/+Hb0//b187/3NjP/9zYz//c2M//3NjP/9zYz//c2M//3NjP/9zYz//Z1cz/49/W/9DMw/83
NTD/EA4L/x8dGf8bGRX/GxkV/xsZFf8bGRX/cT1S/+pvp//eap7/4Gug/+BroP/ga6D/4Gug/+Br
oP/fa5//52uk/3RKWf/JzMD/4dvT/9vXzv/c2M//3NjP/9zYz//c2M//3NjP/9zYz//c2M//3NjP
/9zYz//Z1cz/49/W/9DMw/83NTD/EA4L/x8dGf8bGRX/GxkV/xsZFf9xPVL/6m+n/95qnv/ga6D/
4Gug/+BroP/ga6D/4Gug/99rn//na6T/dEpZ/8nMwP/h29P/29fO/9zYz//c2M//3NjP/9zYz//c
2M//3NjP/9zYz//c2M//3NjP/9zYz//Z1cz/49/W/9DMw/83NTD/EA4L/x8dGf8bGRX/GxkV/3E9
Uv/qb6f/3mqe/+BroP/ga6D/4Gug/+BroP/ga6D/32uf/+drpP90Sln/yczA/+Hb0//b187/3NjP
/9zYz//c2M//3NjP/9zYz//c2M//3NjP/9zYz//c2M//3NjP/9zYz//Z1cz/49/W/9DMw/83NTD/
EA4L/x8dGf8bGRX/cDxR/+hupf/caZ3/3mqe/95qnv/eap7/3mqe/95qnv/cap7/5Wqi/3NKWf/H
yr7/3tnR/9nVzP/a1s3/2tbN/9rWzf/a1s3/2tbN/9rWzf/a1s3/2tbN/9rWzf/a1s3/2tbN/9rW
zf/X08r/4d3U/83Kwf80Mi7/DQsI/x4cF/91P1X/9XSu/+hupf/qb6f/6m+n/+pvp//qb6f/6m+n
/+hvpv/xb6v/eU1d/9LWyf/r5dz/5eHY/+bi2P/m4tj/5uLY/+bi2P/m4tj/5uLY/+bi2P/m4tj/
5uLY/+bi2P/m4tj/5uLY/+bi2f/j39X/7enf/9rWzf9EQj3/FRMP/0ApL/91P1X/cDxR/3E9Uv9x
PVL/cT1S/3E9Uv9xPVL/cD1R/3Q9U/9CLjP/Z2dg/3FuaP9vbGb/b2xm/29sZv9vbGb/b2xm/29s
Zv9vbGb/b2xm/29sZv9vbGb/b2xm/29sZv9vbGb/b2xm/3BtZ/9saWP/eHVv/1xaVP8XFRH/AAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAA=
`
