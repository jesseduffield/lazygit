// Package brainpool implements Brainpool elliptic curves.
// Implementation of rcurves is from github.com/ebfe/brainpool
// Note that these curves are implemented with naive, non-constant time operations
// and are likely not suitable for environments where timing attacks are a concern.
package brainpool

import (
	"crypto/elliptic"
	"math/big"
	"sync"
)

var (
	once                   sync.Once
	p256t1, p384t1, p512t1 *elliptic.CurveParams
	p256r1, p384r1, p512r1 *rcurve
)

func initAll() {
	initP256t1()
	initP384t1()
	initP512t1()
	initP256r1()
	initP384r1()
	initP512r1()
}

func initP256t1() {
	p256t1 = &elliptic.CurveParams{Name: "brainpoolP256t1"}
	p256t1.P, _ = new(big.Int).SetString("A9FB57DBA1EEA9BC3E660A909D838D726E3BF623D52620282013481D1F6E5377", 16)
	p256t1.N, _ = new(big.Int).SetString("A9FB57DBA1EEA9BC3E660A909D838D718C397AA3B561A6F7901E0E82974856A7", 16)
	p256t1.B, _ = new(big.Int).SetString("662C61C430D84EA4FE66A7733D0B76B7BF93EBC4AF2F49256AE58101FEE92B04", 16)
	p256t1.Gx, _ = new(big.Int).SetString("A3E8EB3CC1CFE7B7732213B23A656149AFA142C47AAFBC2B79A191562E1305F4", 16)
	p256t1.Gy, _ = new(big.Int).SetString("2D996C823439C56D7F7B22E14644417E69BCB6DE39D027001DABE8F35B25C9BE", 16)
	p256t1.BitSize = 256
}

func initP256r1() {
	twisted := p256t1
	params := &elliptic.CurveParams{
		Name:    "brainpoolP256r1",
		P:       twisted.P,
		N:       twisted.N,
		BitSize: twisted.BitSize,
	}
	params.Gx, _ = new(big.Int).SetString("8BD2AEB9CB7E57CB2C4B482FFC81B7AFB9DE27E1E3BD23C23A4453BD9ACE3262", 16)
	params.Gy, _ = new(big.Int).SetString("547EF835C3DAC4FD97F8461A14611DC9C27745132DED8E545C1D54C72F046997", 16)
	z, _ := new(big.Int).SetString("3E2D4BD9597B58639AE7AA669CAB9837CF5CF20A2C852D10F655668DFC150EF0", 16)
	p256r1 = newrcurve(twisted, params, z)
}

func initP384t1() {
	p384t1 = &elliptic.CurveParams{Name: "brainpoolP384t1"}
	p384t1.P, _ = new(big.Int).SetString("8CB91E82A3386D280F5D6F7E50E641DF152F7109ED5456B412B1DA197FB71123ACD3A729901D1A71874700133107EC53", 16)
	p384t1.N, _ = new(big.Int).SetString("8CB91E82A3386D280F5D6F7E50E641DF152F7109ED5456B31F166E6CAC0425A7CF3AB6AF6B7FC3103B883202E9046565", 16)
	p384t1.B, _ = new(big.Int).SetString("7F519EADA7BDA81BD826DBA647910F8C4B9346ED8CCDC64E4B1ABD11756DCE1D2074AA263B88805CED70355A33B471EE", 16)
	p384t1.Gx, _ = new(big.Int).SetString("18DE98B02DB9A306F2AFCD7235F72A819B80AB12EBD653172476FECD462AABFFC4FF191B946A5F54D8D0AA2F418808CC", 16)
	p384t1.Gy, _ = new(big.Int).SetString("25AB056962D30651A114AFD2755AD336747F93475B7A1FCA3B88F2B6A208CCFE469408584DC2B2912675BF5B9E582928", 16)
	p384t1.BitSize = 384
}

func initP384r1() {
	twisted := p384t1
	params := &elliptic.CurveParams{
		Name:    "brainpoolP384r1",
		P:       twisted.P,
		N:       twisted.N,
		BitSize: twisted.BitSize,
	}
	params.Gx, _ = new(big.Int).SetString("1D1C64F068CF45FFA2A63A81B7C13F6B8847A3E77EF14FE3DB7FCAFE0CBD10E8E826E03436D646AAEF87B2E247D4AF1E", 16)
	params.Gy, _ = new(big.Int).SetString("8ABE1D7520F9C2A45CB1EB8E95CFD55262B70B29FEEC5864E19C054FF99129280E4646217791811142820341263C5315", 16)
	z, _ := new(big.Int).SetString("41DFE8DD399331F7166A66076734A89CD0D2BCDB7D068E44E1F378F41ECBAE97D2D63DBC87BCCDDCCC5DA39E8589291C", 16)
	p384r1 = newrcurve(twisted, params, z)
}

func initP512t1() {
	p512t1 = &elliptic.CurveParams{Name: "brainpoolP512t1"}
	p512t1.P, _ = new(big.Int).SetString("AADD9DB8DBE9C48B3FD4E6AE33C9FC07CB308DB3B3C9D20ED6639CCA703308717D4D9B009BC66842AECDA12AE6A380E62881FF2F2D82C68528AA6056583A48F3", 16)
	p512t1.N, _ = new(big.Int).SetString("AADD9DB8DBE9C48B3FD4E6AE33C9FC07CB308DB3B3C9D20ED6639CCA70330870553E5C414CA92619418661197FAC10471DB1D381085DDADDB58796829CA90069", 16)
	p512t1.B, _ = new(big.Int).SetString("7CBBBCF9441CFAB76E1890E46884EAE321F70C0BCB4981527897504BEC3E36A62BCDFA2304976540F6450085F2DAE145C22553B465763689180EA2571867423E", 16)
	p512t1.Gx, _ = new(big.Int).SetString("640ECE5C12788717B9C1BA06CBC2A6FEBA85842458C56DDE9DB1758D39C0313D82BA51735CDB3EA499AA77A7D6943A64F7A3F25FE26F06B51BAA2696FA9035DA", 16)
	p512t1.Gy, _ = new(big.Int).SetString("5B534BD595F5AF0FA2C892376C84ACE1BB4E3019B71634C01131159CAE03CEE9D9932184BEEF216BD71DF2DADF86A627306ECFF96DBB8BACE198B61E00F8B332", 16)
	p512t1.BitSize = 512
}

func initP512r1() {
	twisted := p512t1
	params := &elliptic.CurveParams{
		Name:    "brainpoolP512r1",
		P:       twisted.P,
		N:       twisted.N,
		BitSize: twisted.BitSize,
	}
	params.Gx, _ = new(big.Int).SetString("81AEE4BDD82ED9645A21322E9C4C6A9385ED9F70B5D916C1B43B62EEF4D0098EFF3B1F78E2D0D48D50D1687B93B97D5F7C6D5047406A5E688B352209BCB9F822", 16)
	params.Gy, _ = new(big.Int).SetString("7DDE385D566332ECC0EABFA9CF7822FDF209F70024A57B1AA000C55B881F8111B2DCDE494A5F485E5BCA4BD88A2763AED1CA2B2FA8F0540678CD1E0F3AD80892", 16)
	z, _ := new(big.Int).SetString("12EE58E6764838B69782136F0F2D3BA06E27695716054092E60A80BEDB212B64E585D90BCE13761F85C3F1D2A64E3BE8FEA2220F01EBA5EEB0F35DBD29D922AB", 16)
	p512r1 = newrcurve(twisted, params, z)
}

// P256t1 returns a Curve which implements Brainpool P256t1 (see RFC 5639, section 3.4)
func P256t1() elliptic.Curve {
	once.Do(initAll)
	return p256t1
}

// P256r1 returns a Curve which implements Brainpool P256r1 (see RFC 5639, section 3.4)
func P256r1() elliptic.Curve {
	once.Do(initAll)
	return p256r1
}

// P384t1 returns a Curve which implements Brainpool P384t1 (see RFC 5639, section 3.6)
func P384t1() elliptic.Curve {
	once.Do(initAll)
	return p384t1
}

// P384r1 returns a Curve which implements Brainpool P384r1 (see RFC 5639, section 3.6)
func P384r1() elliptic.Curve {
	once.Do(initAll)
	return p384r1
}

// P512t1 returns a Curve which implements Brainpool P512t1 (see RFC 5639, section 3.7)
func P512t1() elliptic.Curve {
	once.Do(initAll)
	return p512t1
}

// P512r1 returns a Curve which implements Brainpool P512r1 (see RFC 5639, section 3.7)
func P512r1() elliptic.Curve {
	once.Do(initAll)
	return p512r1
}
