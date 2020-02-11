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

module.exports = {RemoveByAttribute};