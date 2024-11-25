function UpdateDeviceId(address) {
	
    Swal.fire({
        title: '更改设备标识?',
        text: "请输入新的设备标识",
        input: "text",
	  inputAttributes: {
	    autocapitalize: "off"
	  },
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonText: 'Update',
        cancelButtonText: 'Cancel',
	  preConfirm: async (devid) => {
		let formData = 
		{
			'mac_address': atob(address),
			'devicename': devid
		}
      	const url = '/device';
		const initDetails = {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			  },
			body: JSON.stringify(formData),
			mode: "cors",
		}
		Swal.fire({
               title: 'Processing command...',
               onBeforeOpen: () => {
                   Swal.showLoading()
               }
           	});
		let response = await fetch(url, initDetails);
		Swal.close();
		return response;
	  },
	  allowOutsideClick: () => !Swal.isLoading()
	});
}