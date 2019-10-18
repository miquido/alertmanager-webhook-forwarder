local input = {
  "attachments": [
    {
      "color": null,
      "fields": [
        {
          "short": true,
          "title": "Cluster",
          "value": "youmap-development-main"
        },
        {
          "short": true,
          "title": "Service",
          "value": "youmap-development-api-node"
        },
        {
          "short": true,
          "title": "Tag",
          "value": "0.0.1-dev-119"
        }
      ],
      "pretext": "Deployment has started"
    }
  ],
  "username": "ECS Deploy"
};

local iconsForLabelsAndAnnotations = {
    Cluster: "BOOKMARK",
    Service: "DESCRIPTION",
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

local attachment = input.attachments[0];
local fieldsMap = std.foldl(function(x, y) x { [y.title]: y.value }, attachment.fields, {});

{
    cards: [
        {
            name: input.username,
            header: {
                title: attachment.pretext,
                imageUrl: 'https://miquido.github.io/alertmanager-webhook-forwarder/icons/aws_ecs_icon.png',
            },
            sections: [
                {
                    header: 'Details',
                    widgets: makeWidgets(fieldsMap),
                },
            ]
        }
    ]
}
