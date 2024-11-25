

function KillMe(address) {
    Swal.fire({
        title: 'Are you sure?',
        text: "The device will be disconnected and history will be cleaned.",
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonColor: '#d64130',
        confirmButtonText: 'Kill Me',
        cancelButtonText: 'Cancel'
    }).then((result) => {
        if (result.value) {
            Swal.fire({
                title: 'Processing command...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

            SendCommand(address, "killme")
                .then(response => {
                    if (!response.ok) {
                        throw Error(response.statusText);
                    }
                    return response.text();
                })
                .then(response => {
                    
                    	let formData = 
				{
					'mac_address': atob(address)
				}
				    	const url = '/device';
				const initDetails = {
					method: 'DELETE',
					headers: {
						'Content-Type': 'application/json'
					  },
					body: JSON.stringify(formData),
					mode: "cors",
				}
				let response = await fetch(url, initDetails);
				
			  Swal.close();
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