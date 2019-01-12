// @flow

// global: app templates $ JQuery

export function toHTML(template_name: string, data: Object) {
  if (data == null) {
    data = {};
  }
  data.fmt = app.fmt;
  let template =
    templates[template_name] || templates["component." + template_name];
  if (!template)
    throw new Error(
      `template not found: ${template_name} || component.${template_name}`
    );
  let result = template.render(data, templates);
  delete data.cache;
  delete data.fmt;
  return result;
}

export function Template(
  template_name: string,
  data: Object,
  overrides?: Object
): JQuery {
  if (overrides) {
    data = app.utils.fill(overrides, data);
  }
  return $(toHTML(template_name, data));
}
