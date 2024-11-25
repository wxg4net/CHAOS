function TakeScreenshot(address) {
    Swal.fire({
        title: 'Processing screenshot...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    SendCommand(address, "screenshot")
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
            return response.text();
        })
        .then(response => {
            Swal.close();
            Swal.fire({
		width: '90%',  // 设置宽度
		heightAuto: true, // 启用自定义高度
		html: "<img src=\"/download/"+response+"\" style=\"width: 100%\" onerror=\"this.src='data:image/svg+xml;base64,PHN2ZyBjbGFzcz0iaWNvbiIgc3R5bGU9IndpZHRoOiBhdXRvO2hlaWdodDogYXV0bzt2ZXJ0aWNhbC1hbGlnbjogbWlkZGxlO2ZpbGw6IGN1cnJlbnRDb2xvcjtvdmVyZmxvdzogaGlkZGVuOyIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9IjI1OTUiPjxwYXRoIGQ9Ik01MTEuNjggNjg1LjEyTDQzOC40IDc1OC40YTMyIDMyIDAgMCAxLTIyLjQgOS42IDMyIDMyIDAgMCAxLTIyLjcyLTU0LjcyTDQ2Ni41NiA2NDBsLTczLjI4LTczLjI4YTMyIDMyIDAgMCAxIDQ1LjEyLTQ1LjEybDczLjI4IDczLjI4IDczLjYtNzMuNmEzMiAzMiAwIDEgMSA0NS4xMiA0NS4xMkw1NTYuOCA2NDBsNzMuNiA3My42QTMyIDMyIDAgMCAxIDYwOCA3NjhhMzIgMzIgMCAwIDEtMjIuNzItOS4yOGwtNzMuNi03My42ek03NjggNjcyaC02NGEzMiAzMiAwIDAgMSAwLTY0aDMydi0yODhIMjg4djI4OGgzMmEzMiAzMiAwIDAgMSAwIDY0SDI1NmEzMiAzMiAwIDAgMS0zMi0zMlYyODhhMzIgMzIgMCAwIDEgMzItMzJoNTEyYTMyIDMyIDAgMCAxIDMyIDMydjM1MmEzMiAzMiAwIDAgMS0zMiAzMnoiIHAtaWQ9IjI1OTYiPjwvcGF0aD48L3N2Zz4='\" />",
		showCloseButton: true,
		showCancelButton: false,
		focusConfirm: true,
		});
        }).catch(err => {
        console.log('Error: ', err);
        Swal.fire({
            icon: 'error',
            title: 'Ops...',
            text: 'Error processing screenshot!',
            footer: err
        });
    });
}