package hangouts_chat

import (
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/message_template"
)

var DefaultTemplateJsonnet = message_template.MessageTemplate{
	Type: message_template.Jsonnet,
	Template: `
local cardName = 'Jsonnet Rulez';

{
  cards: [{
    name: cardName,
    header: {
      subtitle: 'Some alerts have occurred!',
      title: 'sometitle',
    },
    sections: [{
      header: std.join(' ', [$.cards[0].header.title, 'Section Header']),
      widgets: [{
        keyValue: { bottomLabel: 'BottomLabel', content: 'Content', icon: 'DOLLAR', topLabel: 'TopLabel' },
      }],
    }],
  }],
  text: std.extVar('input').text,
}
`,
}

var DefaultTemplateGoTemplateYaml = message_template.MessageTemplate{
	Type: message_template.GoTemplateYAML,
	Template: `
cards:
- name: card-name-yaml
  header:
    subtitle: Some alerts have occurred!
    title: AlertManager
  sections:
  - header: Section Header
    widgets:
    - keyValue:
        bottomLabel: BottomLabel
        content: Content
        icon: DOLLAR
        topLabel: TopLabel
text: {{ .Text | toYaml | indent 4 }}
text2:
{{ .Text | indent 8 }}
`,
}

var DefaultTemplateGoTemplateText = message_template.MessageTemplate{
	Type:     message_template.GoTemplateText,
	Template: "{{ .Text }}",
}

var DefaultTemplateAlertmanger = message_template.MessageTemplate{
	Type: message_template.Jsonnet,
	Template: `
local input = std.extVar('input');
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
`,
}
