import React from 'react';
import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import manifest from './manifest';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';

const Icon = () => <i className='icon fa fa-plug'/>;

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        registry.registerChannelHeaderButtonAction(

            // icon - JSX element to use as the button's icon
            <Icon/>,

            // action - a function called when the button is clicked, passed the channel and channel member as arguments
            // null,
            () => {
                alert('Hello World!');
            },

            // dropdown_text - string or JSX element shown for the dropdown button description
            'Hello World',
            'Hello World',
        );
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
