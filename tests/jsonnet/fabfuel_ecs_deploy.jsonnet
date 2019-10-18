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
                imageUrl: "data:image/png;base64,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",
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
