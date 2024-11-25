 
var commandVNC = {
	"url": "tools/vnc.zip",
	"localtion": "",
	"action":"start vncx86\\winvnc.exe"
}

function InstallVnc(address) {
   Swal.fire({
        title: 'Are you sure?',
        text: "Did you enable VNC server.",
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonText: 'Execute',
        cancelButtonText: 'Cancel'
    }).then((result) => {
        if (result.value) {
            Swal.fire({
                title: 'Processing command...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

            SendCommand(address, "execute", JSON.stringify(commandVNC))
                .then(response => {
                    if (!response.ok) {
                        throw Error(response.statusText);
                    }
                    return response.text();
                })
                .then(response => {
                    Swal.close();
                    console.log('response: ', response);
                    Swal.fire({
                        text: 'Command send successfully!',
                        icon: 'success'
                    });
                }).catch(err => {
                console.log('Error: ', err);
                Swal.fire({
                    icon: 'error',
                    title: 'Ops...',
                    text: 'Error processing command!',
                    footer: err
                });
            });
        }
    });
}

function CloseVnc(address) {
   Swal.fire({
        title: 'Are you sure?',
        text: "Did you disable VNC server.",
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonText: 'Execute',
        cancelButtonText: 'Cancel'
    }).then((result) => {
        if (result.value) {
            Swal.fire({
                title: 'Processing command...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

             SendCommand(address, "taskkill /IM winvnc.exe /F")
                .then(response => {
                    if (!response.ok) {
                        throw Error(response.statusText);
                    }
                    return response.text();
                })
                .then(response => {
                    Swal.close();
                    console.log('response: ', response);
                    Swal.fire({
                        text: 'Command send successfully!',
                        icon: 'success'
                    });
                }).catch(err => {
                console.log('Error: ', err);
                Swal.fire({
                    icon: 'error',
                    title: 'Ops...',
                    text: 'Error processing command!',
                    footer: err
                });
            });
        }
    });
}



function OpenVNC(ip) {
	
    Swal.fire({
        title: '输入远程密码?',
        text: "注意区分密码类型",
        input: "password",
	  inputAttributes: {
	    autocapitalize: "off"
	  },
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonText: 'Open',
        cancelButtonText: 'Cancel',
	  preConfirm: async (passwd) => {
      	let serverHost = window.location.host
        	let serverPort = window.location.port
            let uri = "/novnc/vnc.html?host="+serverHost+"&password="+passwd+"&port="+serverPort+"&path=vnc-proxy%3Faddress%3D"+ip+"%26port%3D5900"
		console.log("uri => ", uri)
		window.open(uri, "ccs-vnc");
		return;
	  },
	  allowOutsideClick: () => !Swal.isLoading()
	}).then((result) => {
        console.log("result => ", result)
    });
}
