import { merge } from "lodash";

export const convertBoardApi2Client = data => {
    return {
        fqbn: data.fqbn,
        get config() {
            if (!data.config_options) return {};
            const config = {};
            data.config_options.forEach((con) => {//TODO this is buggy, .value from an empty object is errror-prone.
                config[con.option] = (con.values.find(val => val.selected) || {}).value;
            })
            return config;
        },
        get props() {
            if (!data.properties) return {};
            if (!data.config_options || !data.properties.menu) return data.properties;
            return Object.keys(this.config).reduce((a, i) => merge(
                a,
                ((data.properties.menu[i] || {})[this.config[i]] || {}),
            ), data.properties);
        }
    }
}