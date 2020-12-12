export const valid = () => {
  var result = true
  Array.from(document.querySelectorAll('input'))
    .forEach(i => { if (!i.checkValidity()) result = false })
  return result
}

export const post = (url, data) => {
  let json = {}
  if (data) Object.keys(data).forEach(key => data[key] !== '' && (json[key] = data[key]))
  return fetch(url, {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(json)
  })
}

export const validateEmail = email => {
  //https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
  const re = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  return re.test(email)
}
