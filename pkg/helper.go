package pkg

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomKey(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)[:length]
}

func CreateJsonFile(input JSONRequestSign) (map[string]string, error) {
	jsonData, err := json.MarshalIndent(input.Request.Setup, "", "  ")
	if err != nil {
		return nil, err
	}

	pwd, _ := os.Getwd()
	codeReques := GenerateRandomKey(10)
	path := fmt.Sprintf("%s/public/storage/json/%s.json", pwd, codeReques)
	file, err2 := os.Create(path)
	if err2 != nil {
		return nil, err2
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"code": codeReques,
	}, nil
}

func GetBlockedIPs(maxRequests int) map[string]interface{} {
	blockedIPs := []string{}
	RateLimitMap.Range(func(key, value interface{}) bool {
		data := value.(*struct {
			count     int
			lastReset time.Time
		})
		if data.count >= maxRequests {
			blockedIPs = append(blockedIPs, key.(string))
		}
		return true
	})

	return fiber.Map{"blocked_ips": blockedIPs}
}

func Purchase(ip string) map[string]interface{} {
	if ip == "" {
		return nil
	}

	RateLimitMap.Delete(ip)
	return map[string]interface{}{
		"ip": ip,
	}
}

func RandomIntegerString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := ""
	for i := 0; i < length; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}

func SaveBase64ToFile(base64Str, filename string) error {
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, decoded, 0644)
}

type OCRResponse struct {
	Text  string `json:"text"`
	Error string `json:"error,omitempty"`
}

func RunPythonOCR(imagePath string) (OCRResponse, error) {
	pwd, _ := os.Getwd()
	scriptPath := fmt.Sprintf("%s/pkg/ocr/ocr.py", pwd)
	cmd := exec.Command("python", scriptPath, imagePath)

	// Ambil output dari Python
	output, err := cmd.Output()
	if err != nil {
		return OCRResponse{}, err
	}

	// Parse JSON hasil dari Python
	var result OCRResponse
	err = json.Unmarshal(output, &result)
	if err != nil {
		return OCRResponse{}, err
	}

	return result, nil
}

func SplitLines(text string) []string {
	return strings.Split(text, "\n")
}

func ExtractField(lines []string, keyword string) string {
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			if i+1 < len(lines) {
				return strings.TrimSpace(lines[i+1])
			}
		}
	}
	return ""
}

type ParsedOCRResponse struct {
	Provinsi           string `json:"provinsi"`
	Kabupaten          string `json:"kabupaten"`
	NIK                string `json:"nik"`
	Nama               string `json:"nama"`
	TempatTanggalLahir string `json:"tempat_tanggal_lahir"`
	JenisKelamin       string `json:"jenis_kelamin"`
	GolonganDarah      string `json:"golongan_darah"`
	Alamat             string `json:"alamat"`
	RT_RW              string `json:"rt_rw"`
	KelurahanDesa      string `json:"kelurahan_desa"`
	Kecamatan          string `json:"kecamatan"`
	Agama              string `json:"agama"`
	StatusPerkawinan   string `json:"status_perkawinan"`
	Pekerjaan          string `json:"pekerjaan"`
	Kewarganegaraan    string `json:"kewarganegaraan"`
	BerlakuHingga      string `json:"berlaku_hingga"`
}

func ParseOCRText(ocrText string) ParsedOCRResponse {
	lines := SplitLines(ocrText)
	data := ParsedOCRResponse{}

	// Mapping berdasarkan urutan umum data KTP
	if len(lines) > 0 {
		data.Provinsi = ExtractField(lines, "PROVINSI")
		data.Kabupaten = ExtractField(lines, "KABUPATEN")
		data.NIK = ExtractField(lines, "NIK")
		data.Nama = ExtractField(lines, "Nama")
		data.TempatTanggalLahir = ExtractField(lines, "Tarmpat/TglLahir")
		data.JenisKelamin = ExtractField(lines, "Jenis kelamin")
		data.GolonganDarah = ExtractField(lines, "Gol.Darah")
		data.Alamat = ExtractField(lines, "Alamat")
		data.RT_RW = ExtractField(lines, "RT/RW")
		data.KelurahanDesa = ExtractField(lines, "Kal/Desa")
		data.Kecamatan = ExtractField(lines, "Kecamatan")
		data.Agama = ExtractField(lines, "Agama")
		data.StatusPerkawinan = ExtractField(lines, "Status Perkawinar")
		data.Pekerjaan = ExtractField(lines, "Pekerjaan")
		data.Kewarganegaraan = ExtractField(lines, "Kewarganegaraan")
		data.BerlakuHingga = ExtractField(lines, "Berlaku Hingga")
	}

	return data
}
