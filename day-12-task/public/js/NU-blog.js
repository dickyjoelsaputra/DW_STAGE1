let dataBlog = []

function addBlog(event) {
    event.preventDefault()

    let title = document.getElementById("inputProject").value
    let content = document.getElementById("inputDescription").value
    let image = document.getElementById("inputFile").files[0]
    let startDate = document.getElementById("startDate").value
    let endDate = document.getElementById("endDate").value
    let nodejs = document.getElementById("inlineCheckbox1").checked
    let reactjs = document.getElementById("inlineCheckbox2").checked
    let nextjs = document.getElementById("inlineCheckbox3").checked
    let typescript = document.getElementById("inlineCheckbox4").checked

    // console.log(nodejs)
    // console.log(reactjs)
    // console.log(nextjs)
    // console.log(typescript)

    // if (nodejs == true) {
    //     nodejs = `<i class="fa-brands fa-node-js"></i>`
    // } else if (nodejs = false){
    //     ``
    // }
    
    // if(reactjs == true){
    //     reactjs = `<i class="fa-brands fa-react"></i>`
    // } else if (reactjs = false){
    //     ``
    // }
    
    // if(nextjs == true){
    //     nextjs =`<i class="fa-brands fa-js"></i>`
    // } else if (nextjs = false){
    //     ``
    // }
    
    // if(typescript == true){
    //     typescript = `<i class="fa-brands fa-typo3"></i>`
    // } else if (typescript = false){
    //     ``
    // }


    console.log(nodejs)
    console.log(reactjs)
    console.log(nextjs)
    console.log(typescript)

    // console.log(tech);

    // buat url gambar nantinya tampil
    image = URL.createObjectURL(image)

    
    waktu = getDistanceTime(startDate,endDate)

    let blog = {
        title,
        content,
        image,
        waktu,
        author: "Dicky Joel",
        nodejs,
        reactjs,
        nextjs,
        typescript
        // tech
        
    }

    dataBlog.push(blog)

}

console.log(dataBlog)

function renderBlog() {
    document.getElementById("contents").innerHTML = ''

    for (let i = 0; i < dataBlog.length; i++) {

        document.getElementById("contents").innerHTML += `
        <div class="col-4 my-4 d-flex justify-content-around">
            <div class="card" style="width: 18rem;">
                <img src="${dataBlog[i].image}" class="card-img-top img-fluid"alt="...">
                <div class="card-body">
                <a href="blog-detail.html" class="text-dark text-decoration-none">
                    <h3 class="card-title">${dataBlog[i].title}</h3>
                    <p class="font-weight-light text-muted"> Jangka Waktu : ${dataBlog[i].waktu.bulan} Bulan atau ${dataBlog[i].waktu.hari} Hari</p>
                    <p class="card-text">${dataBlog[i].content}</p>
                </a>
            <h3>
            ${dataBlog[i].nodejs ? `<i class="fa-brands fa-node-js"></i>` : `` }
            ${dataBlog[i].reactjs ? `<i class="fa-brands fa-react"></i>` : `` }
            ${dataBlog[i].nextjs ? `<i class="fa-brands fa-js"></i>` : `` }
            ${dataBlog[i].typescript ? `<i class="fa-brands fa-typo3"></i>` : `` }
            </h3>

                
                <ul class="list-group list-group-flush">
                    
                    <li class="list-group-item">
                        <div class="row">
                            <a href="#" class="btn mx-2 col-5 btn-dark d-inline-block rounded-4">Edit</a>
                            <a href="#" class="btn mx-2 col-5 btn-dark d-inline-block rounded-4">Delete</a>
                        </div>
                    </li>
                    
                </ul>
            
                </div>
            </div>
    </div>
        `
    }
}


// function getFullTime(time) {
//     // time = new Date()
//     // console.log(time)

//     let monthName = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
//     // console.log(monthName[9])

//     // 14
//     let date = time.getDate()
//     console.log(date)

//     // 9
//     let monthIndex = time.getMonth()
//     console.log(monthIndex)

//     // 2022
//     let year = time.getFullYear()
//     console.log(year)

//     let hours = time.getHours()
//     let minutes = time.getMinutes()

//     console.log(hours)

//     if (hours <= 9) {
//         hours = "0" + hours
//     } 
    
//     if (minutes <= 9) {
//         minutes = "0" + minutes
//     }

//     // 14 Oct 2022 09:07 WIB
//     return `${date} ${monthName[monthIndex]} ${year} ${hours}:${minutes} WIB`
// }

function getDistanceTime(s,e) {

    let sDate = new Date(s)
    let eDate = new Date(e)
    //milisecond

    console.log(sDate);
    console.log(eDate);

    let distance = eDate - sDate;

    console.log(distance);

    let milisecondInSecond = 1000 // milisecond
    let secondInMinute = 60
    let minuteInHours = 60 // 1 jam = 24 detik
    let hoursInDay = 24 // 1 hari = 24 jam
    let daysinMonth = 30
    
    let distanceSecond = Math.floor(distance / milisecondInSecond)
    let distanceMinutes = Math.floor(distanceSecond / secondInMinute)
    let distanceHours = Math.floor(distanceMinutes / minuteInHours)
    let distanceDay = Math.floor(distanceHours / hoursInDay)
    let distanceMonth = Math.floor(distanceDay / daysinMonth)

    let waktu = {
        detik : distanceSecond,
        menit : distanceMinutes,
        jam : distanceHours,
        hari : distanceDay,
        bulan: distanceMonth
    };

    console.log(waktu)

    return waktu
    
    // if (distanceDay > 0) {
    //     return `${distanceDay} day(s) ago`
    // } else if (distanceHours > 0) {
    //     return `${distanceHours} hour(s) ago`
    // } else if (distanceMinutes > 0) {
    //     return `${distanceMinutes} minute(s) ago`
    // } else {
    //     return `${distanceSecond} second(s) ago`
    // }



}

// 1#
setInterval(function() {
    renderBlog()
}, 4000)

// 2#
// setInterval(intervalFunction, 3000)

// function intervalFunction() {
//     renderBlog()
// }