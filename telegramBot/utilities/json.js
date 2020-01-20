const errors = require('./errors');

function Parse(data, method) {
    let json;
    try {
        json = JSON.parse(data);
    } catch (e) {
        errors.Log(e, method);
    }
    return json;
}

function RemoveByAttribute(arr, attr, value) {
    let removed = false;
    let i = arr.length;
    while (i--) {
        if (arr[i]
            && arr[i].hasOwnProperty(attr)
            && (arguments.length > 2 && arr[i][attr] === value)) {

            arr.splice(i, 1);
            removed = true;

        }
    }
    return {newList: arr, removed: removed};
}

module.exports = {Parse, RemoveByAttribute};