# Git GUI - UI Specification

## UI Layout & Specifications

### Window Dimensions

**Default Window Size:**
- Width: `1200px`
- Height: `800px`
- Minimum Width: `900px`
- Minimum Height: `600px`
- Resizable: Yes
- Remember last window size/position: Optional but recommended

### Layout Structure

```
┌─────────────────────────────────────────────────────────────────┐
│ HEADER (60px fixed height)                                      │
├──────────────────────┬──────────────────────────────────────────┤
│ FILE LIST            │ DIFF VIEWER                              │
│ (30% width)          │ (70% width)                              │
│ (flex: fills space)  │ (flex: fills space)                      │
│                      │                                          │
│                      │                                          │
├──────────────────────┴──────────────────────────────────────────┤
│ COMMIT PANEL (180px fixed height)                               │
└─────────────────────────────────────────────────────────────────┘
```

### Detailed Component Specs

#### 1. Header (Top Bar)
**Height:** `60px` (fixed)  
**Background:** `#ffffff` (light mode) / `#1e1e1e` (dark mode)  
**Border Bottom:** `1px solid #e0e0e0` / `#333333`  
**Padding:** `0 20px`  
**Display:** `flex`, `align-items: center`, `justify-content: space-between`

**Left Side:**
- **App Title:** "Git Commit Tool"
- Font: `16px`, weight: `600`
- Color: `#333333` / `#ffffff`

**Right Side (Branch Controls):**
- **Branch Dropdown:**
  - Width: `200px`
  - Height: `36px`
  - Border: `1px solid #d0d0d0` / `#404040`
  - Border radius: `4px`
  - Padding: `8px 12px`
  - Font: `14px`
  - Background: `#f5f5f5` / `#2a2a2a`
  - Dropdown icon: Chevron down (8px)
  
- **New Branch Button:**
  - Margin left: `12px`
  - Height: `36px`
  - Padding: `8px 16px`
  - Border: `1px solid #0066cc`
  - Border radius: `4px`
  - Background: `#0066cc`
  - Color: `#ffffff`
  - Font: `14px`, weight: `500`
  - Hover: Background `#0052a3`

#### 2. File List (Left Column)
**Width:** `30%` of content area (minimum `250px`, maximum `400px`)  
**Background:** `#fafafa` / `#252525`  
**Border Right:** `1px solid #e0e0e0` / `#333333`  
**Overflow:** `auto` (vertical scroll)

**Header Section:**
- Height: `40px`
- Padding: `12px 16px`
- Background: `#f0f0f0` / `#2a2a2a`
- Border bottom: `1px solid #e0e0e0` / `#333333`
- Font: `14px`, weight: `600`
- Text: "Changed Files (N)" where N is count

**File List Item:**
- Height: `40px` (fixed per item)
- Padding: `8px 16px`
- Display: `flex`, `align-items: center`
- Gap between elements: `12px`
- Border bottom: `1px solid #e8e8e8` / `#2d2d2d`
- Cursor: `pointer`
- Hover background: `#f0f0f0` / `#2d2d2d`
- Selected background: `#e3f2fd` / `#1a3a52`

**File Item Structure:**
```
[Checkbox (16px)] [Status Badge (60px)] [Filename (flex-1)]
```

**Checkbox:**
- Size: `16px × 16px`
- Border: `1px solid #999999`
- Border radius: `3px`
- Checked background: `#0066cc`
- Checked color: `#ffffff` (checkmark)

**Status Badge:**
- Width: `60px`
- Height: `20px`
- Border radius: `3px`
- Font: `11px`, weight: `600`, uppercase
- Text align: `center`
- Colors by status:
  - Modified: Background `#fff3cd`, Text `#856404`
  - Added: Background `#d4edda`, Text `#155724`
  - Deleted: Background `#f8d7da`, Text `#721c24`
  - Untracked: Background `#d1ecf1`, Text `#0c5460`

**Filename:**
- Font: `13px`, monospace
- Color: `#333333` / `#cccccc`
- Overflow: `hidden`
- Text overflow: `ellipsis`
- White space: `nowrap`

#### 3. Diff Viewer (Right Column)
**Width:** `70%` of content area  
**Background:** `#ffffff` / `#1e1e1e`  
**Overflow:** `auto` (both directions)  
**Padding:** `20px`

**Header (when file selected):**
- Height: `40px`
- Padding: `12px 20px`
- Background: `#f8f8f8` / `#252525`
- Border bottom: `1px solid #e0e0e0` / `#333333`
- Display: `flex`, `justify-content: space-between`, `align-items: center`

**File Path:**
- Font: `14px`, weight: `600`, monospace
- Color: `#333333` / `#ffffff`

**Close Button [×]:**
- Size: `24px × 24px`
- Border: none
- Background: transparent
- Color: `#666666` / `#999999`
- Font size: `20px`
- Cursor: `pointer`
- Hover background: `#e0e0e0` / `#333333`
- Border radius: `3px`

**Diff Content:**
- Font: `13px`, monospace (e.g., 'Consolas', 'Monaco', 'Courier New')
- Line height: `1.5`
- Tab size: `4 spaces`
- White space: `pre` (preserve formatting)

**Diff Line Styling:**
- **Added lines (+):**
  - Background: `#e6ffed` / `#1a3a2a`
  - Color: `#24292e` / `#aff5b4`
  - Prefix: `+` in `#22863a` / `#7ee787`
  
- **Removed lines (-):**
  - Background: `#ffeef0` / `#3a1a2a`
  - Color: `#24292e` / `#ffa198`
  - Prefix: `-` in `#d73a49` / `#ff7b72`
  
- **Context lines:**
  - Background: `transparent`
  - Color: `#586069` / `#8b949e`
  
- **Hunk headers (@@):**
  - Background: `#f6f8fa` / `#21262d`
  - Color: `#0366d6` / `#58a6ff`
  - Font weight: `600`
  - Padding: `4px 8px`
  - Border radius: `3px`

**Empty State (no file selected):**
- Display: `flex`, center aligned
- Color: `#999999`
- Font size: `14px`
- Text: "Select a file to view diff"

#### 4. Commit Panel (Bottom Panel)
**Height:** `180px` (fixed)  
**Background:** `#f8f8f8` / `#252525`  
**Border Top:** `1px solid #e0e0e0` / `#333333`  
**Padding:** `20px`

**Layout:**
```
┌─────────────────────────────────────────────────────────────────┐
│ Label: "Commit Message" (14px, weight: 600, margin-bottom: 8px)│
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ Textarea (height: 80px)                                     │ │
│ │                                                             │ │
│ └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
│                    [Commit]  [Commit & Push]                    │
│                    (margin-top: 12px, float: right)             │
└─────────────────────────────────────────────────────────────────┘
```

**Commit Message Textarea:**
- Width: `100%`
- Height: `80px`
- Padding: `12px`
- Border: `1px solid #d0d0d0` / `#404040`
- Border radius: `4px`
- Font: `14px`, sans-serif
- Resize: `vertical` (allow user to resize)
- Background: `#ffffff` / `#1e1e1e`
- Color: `#333333` / `#cccccc`
- Placeholder: "Enter commit message..."
- Placeholder color: `#999999`

**Button Container:**
- Display: `flex`
- Gap: `12px`
- Justify content: `flex-end`
- Margin top: `12px`

**Commit Button:**
- Height: `40px`
- Padding: `10px 24px`
- Border: `1px solid #0066cc`
- Border radius: `4px`
- Background: `#ffffff` / `#1e1e1e`
- Color: `#0066cc`
- Font: `14px`, weight: `500`
- Cursor: `pointer`
- Hover: Background `#f0f7ff` / `#1a3a52`
- Disabled: Opacity `0.5`, cursor `not-allowed`

**Commit & Push Button:**
- Height: `40px`
- Padding: `10px 24px`
- Border: `1px solid #0066cc`
- Border radius: `4px`
- Background: `#0066cc`
- Color: `#ffffff`
- Font: `14px`, weight: `500`
- Cursor: `pointer`
- Hover: Background `#0052a3`
- Disabled: Opacity `0.5`, cursor `not-allowed`

#### 5. Branch Dropdown (Expanded State)
**Width:** `200px`  
**Max height:** `300px`  
**Background:** `#ffffff` / `#2a2a2a`  
**Border:** `1px solid #d0d0d0` / `#404040`  
**Border radius:** `4px`  
**Box shadow:** `0 2px 8px rgba(0,0,0,0.15)`  
**Overflow:** `auto`

**Branch Item:**
- Height: `36px`
- Padding: `8px 12px`
- Display: `flex`, `align-items: center`
- Gap: `8px`
- Cursor: `pointer`
- Hover background: `#f0f0f0` / `#333333`

**Current Branch Indicator:**
- Icon: Checkmark (12px)
- Color: `#0066cc`
- Shown before branch name

**Branch Name:**
- Font: `14px`, monospace
- Color: `#333333` / `#cccccc`
- Font weight: `600` for current branch, `400` for others

#### 6. New Branch Modal
**Overlay:**
- Background: `rgba(0, 0, 0, 0.5)`
- Position: `fixed`, covers entire window
- Z-index: `1000`
- Display: `flex`, center aligned

**Modal Box:**
- Width: `400px`
- Background: `#ffffff` / `#2a2a2a`
- Border radius: `8px`
- Box shadow: `0 4px 16px rgba(0,0,0,0.2)`
- Padding: `24px`

**Modal Header:**
- Font: `18px`, weight: `600`
- Color: `#333333` / `#ffffff`
- Margin bottom: `20px`
- Text: "Create New Branch"

**Input Field:**
- Width: `100%`
- Height: `40px`
- Padding: `10px 12px`
- Border: `1px solid #d0d0d0` / `#404040`
- Border radius: `4px`
- Font: `14px`
- Background: `#ffffff` / `#1e1e1e`
- Color: `#333333` / `#cccccc`
- Margin bottom: `20px`
- Placeholder: "feature/my-branch"

**Button Container:**
- Display: `flex`
- Gap: `12px`
- Justify content: `flex-end`

**Cancel Button:**
- Height: `40px`
- Padding: `10px 20px`
- Border: `1px solid #d0d0d0` / `#404040`
- Border radius: `4px`
- Background: `transparent`
- Color: `#666666` / `#999999`
- Font: `14px`, weight: `500`
- Cursor: `pointer`

**Create Button:**
- Height: `40px`
- Padding: `10px 20px`
- Border: none
- Border radius: `4px`
- Background: `#0066cc`
- Color: `#ffffff`
- Font: `14px`, weight: `500`
- Cursor: `pointer`
- Hover: Background `#0052a3`

### Color Palette

**Light Mode:**
```
Primary:        #0066cc (Blue)
Primary Hover:  #0052a3 (Darker Blue)
Background:     #ffffff (White)
Surface:        #f8f8f8 (Light Gray)
Border:         #e0e0e0 (Gray)
Text Primary:   #333333 (Dark Gray)
Text Secondary: #666666 (Medium Gray)
Text Muted:     #999999 (Light Gray)

Status Colors:
Success:        #155724 (Green)
Warning:        #856404 (Orange)
Error:          #721c24 (Red)
Info:           #0c5460 (Teal)
```

**Dark Mode:**
```
Primary:        #0066cc (Blue)
Primary Hover:  #0052a3 (Darker Blue)
Background:     #1e1e1e (Very Dark Gray)
Surface:        #252525 (Dark Gray)
Border:         #333333 (Medium Dark Gray)
Text Primary:   #ffffff (White)
Text Secondary: #cccccc (Light Gray)
Text Muted:     #999999 (Medium Gray)

Status Colors:
Success:        #7ee787 (Light Green)
Warning:        #f0ad4e (Orange)
Error:          #ff7b72 (Red)
Info:           #58a6ff (Light Blue)
```

### Typography

**Font Families:**
- UI Text: `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`
- Code/Monospace: `'Consolas', 'Monaco', 'Courier New', monospace`

**Font Sizes:**
- Large: `18px` (Modal headers)
- Medium: `16px` (App title)
- Base: `14px` (Most UI text)
- Small: `13px` (File list, diff content)
- Tiny: `11px` (Status badges)

**Font Weights:**
- Bold: `600`
- Medium: `500`
- Regular: `400`

### Spacing System

Use consistent spacing scale:
- `4px` - Tiny gaps
- `8px` - Small gaps
- `12px` - Medium gaps
- `16px` - Large gaps
- `20px` - Extra large gaps
- `24px` - Section padding

### Responsive Behavior

**Window Resizing:**
- File list width scales proportionally but stays within `250px - 400px`
- Diff viewer takes remaining space
- Minimum window width enforced at `900px` to prevent layout breaking
- At minimum width, file list is `250px`, diff viewer is `650px`

**Scrolling:**
- File list scrolls vertically when content exceeds height
- Diff viewer scrolls both vertically and horizontally for wide diffs
- Commit panel is fixed at bottom (no scroll)
- Header is fixed at top (no scroll)

**Overflow Handling:**
- Long filenames: truncate with ellipsis
- Long branch names: truncate with ellipsis in dropdown
- Long diff lines: horizontal scroll in diff viewer
- Many files: vertical scroll in file list

### Loading & Error States

**Loading Indicator:**
- Display: Spinner (24px) centered in affected area
- Color: `#0066cc`
- Animation: Rotate 360° in 1 second

**Error Toast/Banner:**
- Position: Top center of window
- Width: `auto`, max `600px`
- Height: `auto`
- Padding: `12px 20px`
- Background: `#f8d7da` / `#3a1a1a`
- Border: `1px solid #f5c6cb` / `#721c24`
- Border radius: `4px`
- Color: `#721c24` / `#ffa198`
- Font: `14px`
- Auto-dismiss after `5 seconds`
- Close button: `×` on right side

**Success Toast:**
- Same as error toast but:
- Background: `#d4edda` / `#1a3a2a`
- Border: `#c3e6cb` / `#155724`
- Color: `#155724` / `#7ee787`

### Accessibility

**Keyboard Navigation:**
- Tab order: Branch dropdown → New Branch button → File list → Diff close button → Commit message → Commit → Commit & Push
- Enter key: Submit forms (new branch, commit)
- Escape key: Close modal/dropdown
- Arrow keys: Navigate file list
- Space: Toggle checkbox in file list

**Focus Indicators:**
- Outline: `2px solid #0066cc`
- Outline offset: `2px`
- Applied to all interactive elements

**ARIA Labels:**
- Add appropriate `aria-label` to icon buttons
- Add `role="listbox"` to file list
- Add `role="dialog"` to modal

### Animation/Transitions

**Keep minimal for performance:**
- Button hover: `background-color 0.15s ease`
- Modal open/close: `opacity 0.2s ease`
- Dropdown open/close: `opacity 0.15s ease`
- File selection: `background-color 0.1s ease`

**No animations for:**
- Diff rendering
- File list updates
- Branch switching
