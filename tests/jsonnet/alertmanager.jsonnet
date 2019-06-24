# see pkg/hangouts_chat/message_templates.go
local input = {
    "version": "4",
    "groupKey": "groupkey",
    "status": "firing",
    "receiver": "webhook-reciever",
    "groupLabels": {
        "alertname": "TestAlert",
        "groupLabel": "GroupLabel"
    },
    "commonLabels": {
        "commonLabel": "CommonLabel"
    },
    "commonAnnotations": {
        "commonAnnotation": "CommonAnnotation"
    },
    "externalURL": "https://google.com",
    "alerts": [
        {
            "status": "firing",
            "labels": {
                "alertname": "TestAlert1",
                "severity": "critical"
            },
            "annotations": {
                "message": "TestAlert1 Message",
                "runbook_url": "http://runbook.docs/url"
            },
            "startsAt": "2002-10-02T15:00:00.05Z",
            "endsAt": "2002-11-02T15:00:00.05Z",
            "generatorURL": "https://google.pl"
        },
        {
            "status": "firing",
            "labels": {
                "alertname": "TestAlert2",
                "severity": "warning"
            },
            "annotations": {
                "message": "TestAlert2 Message",
                "runbook_url": "http://runbook.docs2/url"
            },
            "startsAt": "2002-10-02T15:00:00.05Z",
            "endsAt": "2002-11-02T15:00:00.05Z",
            "generatorURL": "https://google.pl"
        }
    ]
};

local alerts = input.alerts;

local iconsForLabelsAndAnnotations = {
    severity: "BOOKMARK",
    message: "DESCRIPTION",
    alertname: "TICKET",
};

local findIconForLabelOrAnnoation(key) = if std.objectHas(iconsForLabelsAndAnnotations, key)
then iconsForLabelsAndAnnotations[key]
else 'STAR';

local makeKVWidget(name, content) = [{
    keyValue: {
        topLabel: name,
        content: content,
        icon: findIconForLabelOrAnnoation(name),
    }
}];

local makeLongWidget(name, content) = [
    {
        keyValue: {
            content: name,
            icon: findIconForLabelOrAnnoation(name),
        }
    },
    {
        textParagraph: {
            text: content,
        }
    }
];

local makeWidgets(resources) = std.flattenArrays([
    if std.length(resources[name]) > 40
    then makeLongWidget(name, resources[name])
    else makeKVWidget(name, resources[name])
    for name in std.objectFields(resources)
]);

{
    cards: [
        {
            name: alert.labels.alertname,
            header: {
                title: alert.labels.alertname + ' (' + alert.labels.severity + ')',
                subtitle: alert.annotations.message,
                imageUrl: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSoaIqZ1iCr1ZGwcJz9W4RVdaIA_AMsyHA6boVH4mEL3bVaRSzT',
            },
            sections: [
                {
                    header: 'Labels',
                    widgets: makeWidgets(alert.labels),
                },
                {
                    header: 'Annotations',
                    widgets: makeWidgets(alert.annotations),
                },
                {
                    widgets: [
                        {
                            keyValue: {
                                topLabel: 'Starts at',
                                content: alert.startsAt,
                                icon: 'CLOCK'
                            },
                        },
                    ],
                },
            ] + (
                if std.objectHas(alert.annotations, 'runbook_url') then [
                    {
                        widgets: [{
                            buttons: [{
                                textButton: {
                                    text: 'Open runnbook',
                                    onClick: {
                                        openLink: {
                                            url: alert.annotations.runbook_url,
                                        },
                                    },
                                },
                            }],
                        }],
                    },
                ] else []
            ),
        }
        for alert in alerts
    ] + [{
        sections: [{
            widgets: [{
                buttons: [{
                    textButton: {
                        text: 'Open alertmanager',
                        onClick: {
                            openLink: {
                                url: input.externalURL,
                            },
                        },
                    },
                }],
            }],
        }],
    }],
}
