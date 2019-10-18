local input = {
  attachments: [
    {
      color: 'good',
      fields: [
        {
          short: true,
          title: 'Cluster',
          value: 'youmap-development-main',
        },
        {
          short: true,
          title: 'Service',
          value: 'youmap-development-api-node',
        },
        {
          short: true,
          title: 'Revision',
          value: 141,
        },
        {
          short: true,
          title: 'Duration',
          value: '0:01:57.735686',
        },
      ],
      pretext: 'Deployment finished successfully',
    },
  ],
  username: 'ECS Deploy',
};

local iconsForLabelsAndAnnotations = {
  Cluster: 'BOOKMARK',
  Service: 'DESCRIPTION',
  Duration: 'CLOCK',
  Tag: 'MAP_PIN',
  Revision: 'MAP_PIN',
};

local findIconForLabelOrAnnoation(key) = if std.objectHas(iconsForLabelsAndAnnotations, key)
then iconsForLabelsAndAnnotations[key]
else 'STAR';

local toString(value) =
  if !std.isString(value)
  then '' + value
  else value;

local makeKVWidget(name, content) = [{
  keyValue: {
    topLabel: name,
    content: toString(content),
    icon: findIconForLabelOrAnnoation(name),
  },
}];

local makeLongWidget(name, content) = [
  {
    keyValue: {
      content: toString(content),
      icon: findIconForLabelOrAnnoation(name),
    },
  },
  {
    textParagraph: {
      text: content,
    },
  },
];

local makeWidgets(resources) = std.flattenArrays([
  if std.isString(resources[name]) && std.length(resources[name]) > 40
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
        imageUrl: 'https://miquido.github.io/alertmanager-webhook-forwarder/icons/aws_ecs.png',
      },
      sections: [
        {
          header: 'Details',
          widgets: makeWidgets(fieldsMap),
        },
      ],
    },
  ],
}
