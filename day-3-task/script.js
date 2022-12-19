function sendEmail() {

    // docID = document.getElementById();

    let name = document.getElementById("name").value
    let email = document.getElementById("email").value
    let phone = document.getElementById("phone").value
    let subject = document.getElementById("subject").value
    let message = document.getElementById("message").value
    
    switch(name && email && phone && subject && message == '') {
        case name:
          return alert("mohon isi nama");
        //   break;
        case email:
            alert("mohon isi email");
        // break;
        case phone:
            return alert("mohon isi phone");
        // break;
        case subject:
            return  alert("mohon isi subject");
        // break;
        case message:
            return alert("mohon isi message");
        // break;
        default:
            alert("Terimakasih telah mengisi form");
      }

    let emailReciever= "dickyjoelsaputra@gmail.com"
    let body = `Hello, my name is ${name} and this is my phone number ${phone}, thank you! , i would like to ${message}`

    let a = document.createElement('a')
    // a.href = `https://mail.google.com/mail/?view=cm&fs=1&to=${emailReciever}&su=${subject} - ${name}&body=${body}`
    a.href = `mailto:${emailReciever}?Subject=${subject}&body=${body}`
    a.click()
}