package entities

import "gorm.io/gorm"

type Prismapxcuatroregistro struct {
	gorm.Model
	Eclq02llEmpresa_04           string
	Eclq02llFpres_04             string
	Eclq02llTiporeg_04           string
	Eclq02llMoneda_04            string
	Eclq02llNumcom_04            string
	Eclq02llNumest_04            string
	Eclq02llNroliq_04            string
	Eclq02llFpag_04              string
	Eclq02llTipoliq_04           string
	Eclq02llCasacta              string
	Eclq02llTipcta               string
	Eclq02llCtabco               string
	Eclq02llCfExentoIva          string
	Eclq02llSigno_04_1           string
	Eclq02llLey25063             string
	Eclq02llSigno_04_2           string
	Eclq02llAliIngbru            string
	Eclq02llDtoCampania          string
	Eclq02llSigno_04_3           string
	Eclq02llIva1DtoCampania      string
	Eclq02llSigno_04_4           string
	Eclq02llRetIngbru2           string
	Eclq02llSigno_04_5           string
	Eclq02llAliIngbru2           string
	Filler1                      string
	Filler2                      string
	Filler3                      string
	Filler4                      string
	Filler5                      string
	Eclq02llTasaPex              string
	Eclq02llCargoXLiq            string
	Eclq02llSigno_04_8           string
	Eclq02llIva1CargoXLiq        string
	Eclq02llSigno_04_9           string
	Eclq02llDealer               string
	Eclq02llImpDbCr              string
	Eclq02llSigno_04_10          string
	Eclq02llCfNoReduceIva        string
	Eclq02llSigno_04_11          string
	Eclq02llPercepIbAgip         string
	Eclq02llSigno_04_12          string
	Eclq02llAlicPercepIbAgip     string
	Eclq02llRetenIbAgip          string
	Eclq02llSigno_04_13          string
	Eclq02llAlicRetenIbAgip      string
	Eclq02llSubtotalRetivaRg3130 string
	Eclq02llSigno_04_14          string
	Eclq02llProvIngbru           string
	Eclq02llAdicPlancuo          string
	Eclq02llSigno_04_15          string
	Eclq02llIva1AdPlancuo        string
	Eclq02llSigno_04_16          string
	Eclq02llAdicOpinter          string
	Eclq02llSigno_04_17          string
	Eclq02llIva1AdOpinter        string
	Eclq02llSigno_04_18          string
	Eclq02llAdicAltacom          string
	Eclq02llSigno_04_19          string
	Eclq02llIva1AdAltacom        string
	Eclq02llSigno_04_20          string
	Eclq02llAdicCupmanu          string
	Eclq02llSigno_04_21          string
	Eclq02llIva1AdCupmanu        string
	Eclq02llSigno_04_22          string
	Eclq02llAdicAltacomBco       string
	Eclq02llSgno_04_23           string
	Eclq02llIva1AdAltacomBco     string
	Eclq02llSigno_04_24          string
	Filler6                      string
	Filler7                      string
	Filler8                      string
	Filler9                      string
	Eclq02llAdicMovpag           string
	Eclq02llSigno_04_27          string
	Eclq02llIva1AdicMovpag       string
	Eclq02llSigno_04_28          string
	Eclq02llRetSellos            string
	Eclq02llSigno_29             string
	Eclq02llProvSellos           string
	Eclq02llRetIngbru3           string
	Eclq02llSigno_04_30          string
	Eclq02llAliIngbru3           string
	Eclq02llRetIngbru4           string
	Eclq02llSigno_04_31          string
	Eclq02llAliIngbru4           string
	Eclq02llRetIngbru5           string
	Eclq02llSigno_04_32          string
	Eclq02llAliIngbru5           string
	Eclq02llRetIngbru6           string
	Eclq02llSigno_04_33          string
	Eclq02llAliIngbru6           string
	Eclq02llFiller_04_10         string
	Eclq02llAster_04             string
	Nombrearchivo                string
	PxDosRegistros               []Prismapxdosregistro `gorm:"foreignkey:PrismapxcuatroregistrosId"`
}
