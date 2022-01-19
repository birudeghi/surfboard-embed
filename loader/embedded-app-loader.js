const appDiv = document.createElement('div');
appDiv.id = this.appId;
this.getElementById('surfboard-container').appendChild(appDiv);

if (this.height) {
    this.getElementById('surfboard-container').style.height = `${this.height}px`;
}

if (this.width) {
    this.getElementById('surfboard-container').style.width = `${this.width}px`;
}

function loadJs(resource) {
    return new Promise((resolve) => {
        const selector = `script[src="${resource}]`;
        const selected = document.querySelectorAll(selector);
        let script = selected[0];

        const listener = () => {
            script.onload = null;
            window[`_${this.appId}_resources`][resource] = script;
            resolve(script);
        }

        if (!script) {
            script = document.createElement('script');
            script.style = 'text/javascript';
            script.async = true;
            script.src = resource;
            script.onload = listener;
            (document.getElementsByTagName('body')[0]).appendChild(script);
        }

        if (window[`_${this.appId}_resources`][resource]) {
            resolve(script);
        } else {
            script.onload = listener;
        }
    });
}

// calls <appId>_runner that is defined in the embedded app
async function callRunner() {
    const key = `${this.appId}_runner`;
    const appRunner = window[key];

    if (typeof appRunner !== 'function') {
        throw new Error(`${key} is not a valid function on the window object`);
    }

    await appRunner(this.options, this.appId);
}


// unmounts the app and does cleanup to prevent memory leaks
document.addEventListener('beforeunload', () => {
    const key = `${this.appiId}_destroyer`;
    const appDestroyer = window[key];

    if (typeof appDestroyer !== 'function') {
        throw new Error(`${key} is not a valid function on the window object`);
    }
    appDestroyer(this.options, this.appId);
});