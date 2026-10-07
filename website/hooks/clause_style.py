"""MkDocs hook that reshapes the manual's rendered pages for the site.

The manual stays plain Markdown; the classes added here are what
assets/extra.css styles. Each clause is wrapped in a level-tagged container and
its heading is split into a level badge, the clause id and the summary.
"""

import re

from bs4 import BeautifulSoup
from bs4.element import NavigableString, Tag

HEADING_PATTERN = re.compile(r"^\s*【(MUST|SHOULD|MAY)】\s*(\S+)\s+(.*)$", re.S)
TOC_PATTERN = re.compile(r"^\s*【(MUST|SHOULD|MAY)】\s*(\S+)\s+(.*)$", re.S)
HEADING_TAGS = {"h1", "h2", "h3"}
EXAMPLE_CLASSES = ("gdm-good", "gdm-bad")


def _shorten_toc(items):
    """The heading carries a level badge and a clause id the sidebar does not need."""
    for item in items:
        match = TOC_PATTERN.match(item.title or "")
        if match:
            item.title = match.group(3)
        _shorten_toc(item.children)


def _is_index(page):
    # The i18n plugin keeps the locale folder in src_uri, so is_homepage is not
    # enough to recognise the manual's index page.
    return page.is_homepage or page.file.src_uri.endswith("index.md")


def _tag_clause_heading(soup, heading):
    """Split 【MUST】SEC-001 summary into a badge, an id and the summary text."""
    text = heading.find(string=True)
    if text is None:
        return None
    match = HEADING_PATTERN.match(str(text))
    if match is None:
        return None

    level, identifier, summary = match.groups()
    badge = soup.new_tag("span", attrs={"class": ["gdm-level", "gdm-level--" + level.lower()]})
    badge.string = level
    clause_id = soup.new_tag("span", attrs={"class": "gdm-id"})
    clause_id.string = identifier
    text.replace_with(badge, clause_id, NavigableString(" " + summary))
    return level.lower()


def on_page_context(context, page, config, nav):
    _shorten_toc(page.toc.items)
    if _is_index(page):
        page.meta["hide"] = ["toc"]
    return context


def on_page_content(html, page, config, files):
    soup = BeautifulSoup(html, "html.parser")

    for heading in soup.find_all("h3"):
        level = _tag_clause_heading(soup, heading)
        if level is None:
            continue

        container = soup.new_tag("div")
        container["class"] = ["gdm-clause", "gdm-clause--" + level]
        heading.insert_before(container)

        body = []
        for node in heading.next_siblings:
            if isinstance(node, Tag) and node.name in HEADING_TAGS:
                break
            body.append(node)

        container.append(heading.extract())
        for node in body:
            container.append(node.extract())

        for class_name, block in zip(EXAMPLE_CLASSES, container.select("div.highlight")):
            block["class"].append(class_name)

    if _is_index(page):
        index = soup.new_tag("div")
        index["class"] = ["gdm-index"]
        for node in list(soup.contents):
            index.append(node.extract())
        soup.append(index)

    return str(soup)
