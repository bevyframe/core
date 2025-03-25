const Page = (prop, childs) => {
    document.title = prop.title;
    document.body.className = `body_${prop.themeColor}`;
    for (let element of childs) {
        renderWidget(element);
    }
}
const Navbar = (childs) => ["nav", {"class": "Navbar"}, [...childs]];
const NavItem = (icon, link, alt, active = false) => {
    return ["a", {"class": active? "active" : "inactive", href: link}, [
        ["button", {}, [
            ["span", {"class": "material-symbols-rounded", alt: alt}, [icon]]
        ]]
    ]]
};
const Root = (prop, childs) => ["div", {id: "root"}, [...childs]];
const Margin = (prop) => [];
const Size = {
    pixel: (int) => `${int}px`,
    max_content: "max-content",
};
const Align = {
    center: "center",
};
const Container = (prop, childs) => ["div", prop, [...childs]];
const Button = (type, prop) => ["button", {"class": `button ${type}`, ...prop}, ["DEMO"]];
const Box = (prop, childs) => ["div", {"class": "the_box", ...prop}, [...childs]];
const Line = (childs) => ["p", {}, [...childs]];
const Textbox = (name, prop) => ["input", {name: name, id: name, "class": "textbox", ...prop}, []];