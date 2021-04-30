export interface PluginRegistry {
    registerPostTypeComponent(typeName: string, component: React.ElementType);
    registerChannelHeaderButtonAction(icon, action, dropdownText, tooltipText);

    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
