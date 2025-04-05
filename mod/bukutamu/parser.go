package bukutamu

import (
	"regexp"
	"strings"
)

func ParsePesanFlexible(pesan string) Tamu {
	var tamu Tamu

	lines := strings.Split(pesan, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Asal Tamu / Nama Mitra:"):
			tamu.Kategori = strings.TrimSpace(strings.TrimPrefix(line, "Asal Tamu / Nama Mitra:"))
		case strings.HasPrefix(line, "Alamat yang dituju:"):
			tamu.Tujuan = strings.TrimSpace(strings.TrimPrefix(line, "Alamat yang dituju:"))
		case strings.HasPrefix(line, "Nomor Plat Kendaraan:"):
			tamu.Kendaraan = strings.TrimSpace(strings.TrimPrefix(line, "Nomor Plat Kendaraan:"))
		}
	}
	tamu.BlokRumah = extractBlokRumah(tamu.Tujuan)
	return tamu
}

func extractBlokRumah(alamat string) string {
	alamat = strings.ToUpper(strings.TrimSpace(alamat))

	// Pola regex diperluas untuk menangkap angka + huruf (3B, 12A, dsb)
	patterns := []string{
		`BLO[KC]\\s*([A-G])\s*[\.,\-]?\s*(?:NO(?:\.|MOR)?\s*)?(\d+[A-Z]?)`, // Blok C No. 3B
		`([A-G])[\s\-]?(\d+[A-Z]?)`,                                        // C 10, D-15, E9A
		`RUMAH\s+([A-G])(\d+[A-Z]?)`,                                       // Rumah C12B
		`BLO[KC]([A-G])(\d+[A-Z]?)`,                                        // blokC12
	}

	for _, pat := range patterns {
		re := regexp.MustCompile(pat)
		match := re.FindStringSubmatch(alamat)

		if len(match) >= 3 {
			blok := strings.ToUpper(match[1])
			nomor := strings.ToUpper(match[2])

			if isValidBlok(blok) && isValidNomor(nomor) {
				return blok + nomor
			}
		}
	}

	return "" // Tidak valid
}

func isValidBlok(blok string) bool {
	return strings.Contains("ABCDEFG", blok)
}

func isValidNomor(nomor string) bool {
	re := regexp.MustCompile(`^\d{1,3}[A-Z]?$`) // 1, 12, 3B, 99A, dst
	return re.MatchString(nomor)
}
