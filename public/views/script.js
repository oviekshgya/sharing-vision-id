document.addEventListener("DOMContentLoaded", function () {
    const video = document.getElementById("video");
    const canvas = document.getElementById("canvas");
    const resultDiv = document.getElementById("result");
    const checkAgainBtn = document.getElementById("checkAgainBtn");
    const loadingSpinner = document.querySelector(".loading-spinner");

    let detectionActive = true;
    let isProcessing = false; // Cegah request berulang sebelum hasil didapat

    // Akses Kamera
    navigator.mediaDevices.getUserMedia({ video: true })
        .then(stream => video.srcObject = stream)
        .catch(error => console.error("Gagal mengakses kamera:", error));

    async function captureAndDetect() {
        if (!detectionActive || isProcessing) return; // Stop jika sudah terdeteksi atau masih memproses

        isProcessing = true; // Tandai bahwa request sedang berjalan
        loadingSpinner.style.display = "block"; // Tampilkan spinner

        const context = canvas.getContext("2d");
        context.drawImage(video, 0, 0, canvas.width, canvas.height);

        // Konversi ke Base64
        const imageBase64 = canvas.toDataURL("image/jpeg");

        const username = "Shagya";
        const password = "ShagyaTech";
        const basicAuth = btoa(`${username}:${password}`);

        try {
            const response = await fetch("/demo/detect-live", {
                method: "POST",
                headers: { "Content-Type": "application/json", "Authorization": `Basic ${basicAuth}`,"X-API-KEY": "xzDcUsxhstdalZtbdMz0" },
                body: JSON.stringify({ image: imageBase64 })
            });

            const data = await response.json();

            if (data.status === "success") {
                resultDiv.className = "alert alert-success text-center";
                resultDiv.innerHTML = `<i class="fas fa-check-circle"></i> ${data.message} (${data.recognized.join(", ")})`;
                detectionActive = false;
                showCheckAgainButton();
            } else {
                resultDiv.className = "alert alert-danger text-center";
                resultDiv.innerHTML = `<i class="fas fa-times-circle"></i> ${data.message}`;
            }
        } catch (error) {
            console.error("Error:", error);
        } finally {
            isProcessing = false; // Reset agar bisa mengirim request berikutnya
            loadingSpinner.style.display = "none"; // Sembunyikan spinner
        }
    }

    function showCheckAgainButton() {
        checkAgainBtn.classList.remove("d-none");
        checkAgainBtn.onclick = function () {
            detectionActive = true;
            resultDiv.className = "alert alert-secondary text-center";
            resultDiv.innerHTML = "<p><i class='fas fa-hourglass-half'></i> Menunggu deteksi wajah...</p>";
            checkAgainBtn.classList.add("d-none"); // Sembunyikan tombol setelah klik
        };
    }

    // Jalankan live detection setiap 2 detik (hanya jika tidak sedang diproses)
    setInterval(captureAndDetect, 2000);
});
