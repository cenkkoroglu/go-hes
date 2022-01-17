$(function () {
    $("#checkHes").submit(function (e) {
        e.preventDefault();

        const form = $(this);
        const hesCode = form.find('input[name="hesCode"]').val();

        if (hesCode == "") {
            Swal.fire({
                icon: 'error',
                title: 'Hata...',
                text: 'Lütfen geçerli bir HES kodu girin!',
                heightAuto: false
            })
            return
        } else {
            $.ajax({
                url: "/checkHesCode?hes=" + hesCode,
                type: 'post',
            }).done(function (result) {
                if (result.status == 200) {
                    if (result.data.error !== undefined && result.data.error) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Hata...',
                            text: result.data.error,
                            heightAuto: false
                        })
                    } else {
                        Swal.fire({
                            icon: 'success',
                            html: `<table class="table table-bordered">
                                <tr>
                                    <td>İsim</td>
                                    <td>` + result.data.masked_firstname + `</td>
                                </tr>
                                <tr>
                                    <td>Soyisim</td>
                                    <td>` + result.data.masked_lastname + `</td>
                                </tr>
                                <tr>
                                    <td>TC Numarası</td>
                                    <td>` + result.data.masked_identity_number + `</td>
                                </tr>
                                <tr>
                                    <td>Sağlık Durumu</td>
                                    <td>` + result.data.current_health_status + `</td>
                                </tr>
                                <tr>
                                    <td>Geçerlilik Tarihi</td>
                                    <td>` + result.data.expiration_date + `</td>
                                </tr>
                            </table>`,
                            heightAuto: false
                        })
                    }
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Hata...',
                        text: result.data.title,
                        heightAuto: false
                    })
                }
            });
        }
    });

    $("#confirmationForm").submit(function (e) {
        e.preventDefault();

        const form = $(this);
        const phoneNumber = form.find('input[name="phoneNumber"]').val();

        if (phoneNumber.length != 12) {
            Swal.fire({
                icon: 'error',
                title: 'Hata...',
                text: 'Telefon numarası 905xxxxxxxxx formatında olmalıdır!',
                heightAuto: false
            })
            return
        }

        $.ajax({
            url: "/sendLoginCode?phoneNumber=" + phoneNumber,
            type: 'post',
        }).done(function (result) {
            if (result.status == 201) {
                Swal.fire({
                    title: 'Doğrulama',
                    html: `<input type="text" id="login" class="swal2-input" placeholder="xxxxx">`,
                    confirmButtonText: 'Doğrula',
                    focusConfirm: false,
                    preConfirm: () => {
                        const login = Swal.getPopup().querySelector('#login').value
                        if (!login || login.length != 5) {
                            Swal.showValidationMessage(`Lütfen 5 haneli onay kodunu girin!`)
                        }
                        return login
                    },
                    heightAuto: false
                }).then((result) => {
                    if (result.isDismissed) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Hata...',
                            text: 'Doğrulama iptal edildi!',
                            heightAuto: false
                        })
                    } else {
                        $.ajax({
                            url: "/authenticate?phoneNumber=" + phoneNumber + "&loginCode=" + result.value,
                            type: 'post',
                        }).done(function (result) {
                            if (result.status == 200) {
                                location.reload()
                            } else {
                                Swal.fire({
                                    icon: 'error',
                                    title: 'Hata...',
                                    text: result.data.title,
                                    heightAuto: false
                                })
                            }
                        });
                    }
                })
            } else {
                Swal.fire({
                    icon: 'error',
                    title: 'Hata...',
                    text: result.data.title,
                    heightAuto: false
                })
            }
        });
    });
});
