import axios from "axios";
export function Login(user) {
    axios({
        method: 'POST',
        url: '/login',
        data: user
    }).then(response => {
        console.log(response)
    }).catch(err => {
        console.log(err)
    })
}