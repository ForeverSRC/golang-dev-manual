"""MkDocs hook that wraps each clause in a level-tagged container.

The manual stays plain Markdown; the container and example classes added here
are what assets/extra.css styles.
"""

import re

from bs4 import BeautifulSoup
from bs4.element import Tag

LEVEL_PATTERN = re.compile(r"^【(MUST|SHOULD|MAY)】")
HEADING_TAGS = {"h1", "h2", "h3"}
EXAMPLE_CLASSES = ("gdm-good", "gdm-bad")


def on_page_content(html, page, config, files):
    soup = BeautifulSoup(html, "html.parser")

    for heading in soup.find_all("h3"):
        match = LEVEL_PATTERN.match(heading.get_text(strip=True))
        if match is None:
            continue

        container = soup.new_tag("div")
        container["class"] = ["gdm-clause", "gdm-clause--" + match.group(1).lower()]
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

    return str(soup)
