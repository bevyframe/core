const Size = {
    pixel: (int) => `${int}px`,
    maxContent: "max-content",
};

const Align = {
    center: "center",
};

const Margin = (prop) => { return {
    left: prop.left,
    right: prop.right,
    top: prop.top,
    bottom: prop.bottom,
} };

const Padding = (prop) => { return {
    left: prop.left,
    right: prop.right,
    top: prop.top,
    bottom: prop.bottom,
} };

const loadMargin = (prop) => {
    if ((typeof prop) === "string")
        return `margin: ${prop};`;
    let style = "";
    if (prop.left)
        style += `margin-left: ${prop.left};`;
    if (prop.right)
        style += `margin-right: ${prop.right};`;
    if (prop.top)
        style += `margin-top: ${prop.top};`;
    if (prop.bottom)
        style += `margin-bottom: ${prop.bottom};`;
    return style;
};

const loadPadding = (prop) => {
    if ((typeof prop) === "string")
        return `padding: ${prop};`;
    let style = "";
    if (prop.left)
        style += `padding-left: ${prop.left};`;
    if (prop.right)
        style += `padding-right: ${prop.right};`;
    if (prop.top)
        style += `padding-top: ${prop.top};`;
    if (prop.bottom)
        style += `padding-bottom: ${prop.bottom};`;
    return style;
};

const buildStyle = (prop) => {
    let style = "";
    if (prop.margin)
        style += loadMargin(prop.margin);
    if (prop.padding)
        style += loadPadding(prop.padding);
    if (prop.height)
        style += "height: " + prop.width + ";";
    if (prop.width)
        style += "width: " + prop.width + ";";
    if (prop.textAlign)
        style += "text-align: " + prop.textAlign + ";";
    if (prop.backgroundColor)
        style += "background-color: " + prop.backgroundColor + ";";
    return style;
}

const buildProp = (prop) => {
    let style = buildStyle(prop);
    if (prop.style)
        style = prop.style + ";" + style;
    let eProp = {style: style};
    if (prop.id)
        eProp.id = prop.id;
    if (prop.class)
        eProp.class = prop.class;
    if (prop.onClick)
        eProp.onclick = prop.onClick;
    if (prop.placeholder)
        eProp.placeholder = prop.placeholder;
    if (prop.value)
        eProp.value = prop.value;
    if (prop.type)
        eProp.type = prop.type;
    return eProp;
}

   /* --------------------------------------------------------------------- */
  /* --------------------------------------------------------------------- */
 /* --------------------------------------------------------------------- */

const Page = (prop, childs) => {
    document.title = prop.title;
    document.body.className = `body_${prop.color}`;
    for (let element of childs)
        renderWidget(element);
}

const Navbar = (childs) =>
    ["nav", {"class": "Navbar"}, [...childs]];

const NavItem = (icon, link, alt, active = false) => {
    return ["a", {"class": active? "active" : "inactive", href: link}, [
        ["button", {}, [
            ["span", {"class": "material-symbols-rounded", alt: alt}, [icon]]
        ]]
    ]]
};

const Root = (prop, childs) =>
    ["div", {id: "root", ...buildProp(prop)}, [...childs]];

const Button = (prop) =>
    ["button", {"class": "button", ...buildProp(prop)}, [prop.innerText]];

const SmallButton = (prop) =>
    ["button", {"class": "button small", ...buildProp(prop)}, [prop.innerText]];

const MiniButton = (prop) =>
    ["button", {"class": "button mini", ...buildProp(prop)}, [prop.innerText]];

const Container = (prop, childs) =>
    ["div", buildProp(prop), [...childs]];

const Line = (childs) =>
    ["p", {style: "width: max-content;"}, [...childs]];

const Textbox = (name, prop) =>
    ["input", {name: name, id: name, "class": "textbox", ...buildProp(prop)}, []];

const Title = (text, prop) =>
    ["h1", (prop)?{...buildProp(prop)}:{}, [text]];

  /* --------------------------------------------------------------------- */
 /* --------------------------------------------------------------------- */
/* --------------------------------------------------------------------- */

const Box = (prop, childs) =>
    Container({"class": "the_box", ...prop}, [...childs]);