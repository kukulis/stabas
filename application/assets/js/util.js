function clearTag(tag) {
    while (tag.firstChild) {
        tag.removeChild(tag.firstChild);
    }
}

/**
 *
 * @type map {Map}
 */
function flipMap(map) {
    let resultMap = new Map();
    map.forEach((key, value) => {
        resultMap.set(value, key)
    });

    return resultMap;
}

/**
 *
 * @param dateStr {Date|null}
 */
function parseDate(dateStr) {
    if (dateStr === null) {
        return null;
    }

    return new Date(dateStr)
}

/**
 *
 * @param date {Date|null}
 */
function formatDate(date) {
    if (date === null) {
        return '';
    }

    return date.toISOString()
}
