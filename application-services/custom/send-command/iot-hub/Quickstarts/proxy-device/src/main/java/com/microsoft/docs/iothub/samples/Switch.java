package com.microsoft.docs.iothub.samples;

import com.google.api.client.json.GenericJson;
import com.google.api.client.util.Key;

public class Switch extends GenericJson {
    @Key
    private String SwitchButton;

    public String getSwitchButton() {
        return SwitchButton;
    }

    public void setSwitchButton(String s) {
        this.SwitchButton = s;
    }

    @Override
    public String toString() {
        return "SwitchButton is " + SwitchButton;
    }
}
